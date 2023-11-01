package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"redirect/model"
	"redirect/utils"
)

// FindAllAccessLogs 查询所有日志
func FindAllAccessLogs(ruleId int64, search string, start, end string, page, size int) ([]model.AccessLog, error) {
	var (
		found  []model.AccessLog
		offset = (page - 1) * size
	)
	query := `SELECT * FROM public.access_logs l WHERE 1=1 `
	where := logWhereCondition(ruleId, search, start, end)
	query += where
	query += ` ORDER BY l.id DESC LIMIT $1 OFFSET $2`
	return found, DbSelect(query, &found, size, offset)
}

// FindAllAccessPageTotal 日志分页统计
func FindAllAccessPageTotal(ruleId int64, search, start, end string) (int, error) {
	var count int
	query := `SELECT count(l.id) as total_count FROM public.access_logs l WHERE 1=1 `
	where := logWhereCondition(ruleId, search, start, end)
	query += where
	return count, DbGet(query, &count)
}

func logWhereCondition(ruleId int64, search, start, end string) string {
	var query string
	if ruleId > 0 {
		query += fmt.Sprintf(` AND l.rule_id = %d`, ruleId)
	}
	if !utils.EmptyString(search) {
		likeSearch := "%" + search + "%"
		query += fmt.Sprintf(` AND  l.from_domain LIKE '%s' OR l.to_domain LIKE '%s' OR l.ip LIKE '%s'`,
			likeSearch, likeSearch, likeSearch)
	}
	if !utils.EmptyString(start) {
		query += fmt.Sprintf(` AND l.access_time >= to_date('%s','YYYY-MM-DD')`, start)
	}
	if !utils.EmptyString(end) {
		query += fmt.Sprintf(` AND l.access_time < (to_date('%s','YYYY-MM-DD') + interval  '1day') `, end)
	}
	return query
}

//DeleteAccessLogs 删除日志
func DeleteAccessLogs(ids []int) error {
	sql := `DELETE FROM public.access_logs  WHERE  id in (?)`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return err
	}
	formatQuery := RebindQuery(query)
	return DbExec(formatQuery, args...)
}

//DeleteAccessLogsHistory 删除多少天前历史数据
func DeleteAccessLogsHistory(beforeTime string) error {
	query := `DELETE FROM public.access_logs  WHERE  access_time < $1`
	return DbExec(query, beforeTime)
}

//BatchInsertAccessLogs 批量同步数据
func BatchInsertAccessLogs(logs []model.AccessLog, batchSize int) error {
	db := dbService.Connection
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	query := `INSERT INTO public.access_logs (log_uuid,rule_id,from_domain,to_domain, access_time, ip, user_agent,referer,uv_cookie)`
	query += ` VALUES(:log_uuid,:rule_id,:from_domain,:to_domain,:access_time,:ip,:user_agent,:referer,:uv_cookie)`
	for i := 0; i < len(logs); i += batchSize {
		end := i + batchSize
		if end > len(logs) {
			end = len(logs)
		}
		batch := logs[i:end]
		_, err = tx.NamedExec(query, batch)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

package storage

import (
	"encoding/json"
	"fmt"
	"redirect/model"
	"redirect/utils"
	"strconv"
	"time"
)

const PVKeyPrefix = "pv_stats"
const UVKeyPrefix = "uv_stats"
const IPKeyPrefix = "ip_stats"
const RuleLastPrefix = "rule_rtd"

// InsertStatsMinutes 插入数据
func InsertStatsMinutes(data []model.StatsMinutes) error {
	query := `INSERT INTO public.stats_minutes(access_time,rule_id,pv_num,ip_num,uv_num) ` +
		`VALUES(:access_time,:rule_id,:pv_num,:ip_num,:uv_num)`
	return DbNamedExec(query, data)
}

// DeleteStatsMinutes 删除3天前数据,要统计昨天数据，所以保留3天
func DeleteStatsMinutes() error {
	var yesterday = time.Now().AddDate(0, 0, -3).Format("2006-01-02 15:04:05")
	query := `DELETE FROM  public.stats_minutes WHERE  access_time < $1`
	return DbExec(query, yesterday)
}

// InsertStatsDays 插入数据
func InsertStatsDays(data model.StatsDays) error {
	query := `INSERT INTO public.stats_days(access_time,rule_id,pv_num,ip_num,uv_num) ` +
		`VALUES(:access_time,:rule_id,:pv_num,:ip_num,:uv_num)`
	return DbNamedExec(query, data)
}

// DeleteStatsDays 删除3月前数据
func DeleteStatsDays() error {
	var threeMonth = time.Now().AddDate(0, -3, 0).Format("2006-01-02 15:04:05")
	query := `DELETE FROM  public.stats_days WHERE  access_time < $1`
	return DbExec(query, threeMonth)
}

//DeleteYesterdaysByRuleId 清空昨天数据
func DeleteYesterdaysByRuleId(ruleId int64) error {
	var yesterday = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	query := `DELETE FROM  public.stats_days WHERE rule_id = $1 AND date(access_time) = $2`
	return DbExec(query, ruleId, yesterday)
}

//GetYesterdaysByRuleId 查询昨天数据
func GetYesterdaysByRuleId(ruleId int64) (model.StatsDays, error) {
	var totals model.StatsDays
	var yesterday = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	query := `SELECT * FROM  public.stats_days WHERE rule_id = $1 AND date(access_time) = $2`
	return totals, DbGet(query, &totals, ruleId, yesterday)
}

//TotalTodayNum 统计某天总数
func TotalTodayNum(ruleId int64, date string) (model.StatsDays, error) {
	pvNum := GetRulesPV(ruleId, date)
	uvNum := GetRulesUV(ruleId, date)
	ipNum := GetRulesIP(ruleId, date)
	totals := model.StatsDays{
		PvNum: pvNum,
		UvNum: uvNum,
		IpNum: ipNum,
	}
	return totals, nil
}

//GetHoursChart 获取某天小时统计数据
func GetHoursChart(date string, ruleId int64) ([]model.TotalStats, error) {
	var totals []model.TotalStats
	query := `SELECT to_char(sm.access_time, 'YYYY-MM-DD HH24') as access_time,
	           SUM(sm.pv_num) AS pv_num,
	           SUM(sm.ip_num) AS ip_num,
	           SUM(sm.uv_num) AS uv_num
			FROM public.stats_minutes sm
			WHERE date(sm.access_time) = $1 AND sm.rule_id = $2
			GROUP BY to_char(sm.access_time, 'YYYY-MM-DD HH24')
	       ORDER BY access_time ASC`

	return totals, DbSelect(query, &totals, date, ruleId)
}

//GetDaysChart 获取前几天至今数据
func GetDaysChart(startTime string, ruleId int64) ([]model.TotalStats, error) {
	var totals []model.TotalStats
	query := `SELECT
                to_char(sd.access_time::DATE, 'YYYY-MM-DD') as access_time,
				SUM(sd.pv_num) AS pv_num,
				SUM(sd.ip_num) AS ip_num,
				SUM(sd.uv_num) AS uv_num
			FROM public.stats_days sd
			WHERE sd.access_time >= $1 AND sd.rule_id = $2
			GROUP BY to_char(sd.access_time::DATE, 'YYYY-MM-DD')
            ORDER BY access_time ASC`
	return totals, DbSelect(query, &totals, startTime, ruleId)
}

//GetTopRank 获取排名数据
func GetTopRank(ruleId int64, top int) ([]model.TopRank, error) {
	var totals []model.TopRank
	query := `SELECT
				a.from_domain,
				count(a.id) as num
			FROM public.access_logs  a
			WHERE date(a.access_time) = date(NOW()) AND a.rule_id = $1
            GROUP BY a.from_domain
            ORDER BY num DESC
            LIMIT $2`
	return totals, DbSelect(query, &totals, ruleId, top)
}

//SaveRulesPV 统计pv
func SaveRulesPV(ruleId int64) error {
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", PVKeyPrefix, today, ruleId)
	return RedisIncr(key)
}

// GetRulesPV 获取pv
func GetRulesPV(ruleId int64, date string) int64 {
	var num int64
	//today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", PVKeyPrefix, date, ruleId)

	numStr, _ := RedisGetString(key)
	if numStr == "" {
		num = 0
	}
	numInt, _ := strconv.Atoi(numStr)
	num = int64(numInt)
	return num
}

// SaveRulesUV 利用HyperLogLog保存uv
func SaveRulesUV(ruleId int64, fromDomain string, uvCookie string) error {
	if uvCookie == "" {
		return nil
	}
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", UVKeyPrefix, today, ruleId)
	value := utils.HashPassword(fromDomain+uvCookie, "")
	err := RedisPFAdd(key, value)
	if err != nil {
		return err
	}
	return nil
}

//GetRulesUV 获取uv
func GetRulesUV(ruleId int64, date string) int64 {
	var num int64
	//today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", UVKeyPrefix, date, ruleId)
	num, _ = RedisPFCount(key)
	return num
}

//SaveRulesIP 利用HyperLogLog统计Ip
func SaveRulesIP(ruleId int64, fromDomain string, ip string) error {
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", IPKeyPrefix, today, ruleId)
	value := utils.HashPassword(fromDomain+ip, "")
	err := RedisPFAdd(key, value)
	if err != nil {
		return err
	}
	return nil
}

//GetRulesIP 获取ip
func GetRulesIP(ruleId int64, date string) int64 {
	var num int64
	//today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", IPKeyPrefix, date, ruleId)
	num, _ = RedisPFCount(key)
	return num
}

//SaveTodayLastRecord 保存上一次数据
func SaveTodayLastRecord(ruleId int64, data model.StatsDays) error {
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", RuleLastPrefix, today, ruleId)
	record, _ := json.Marshal(data)
	return RedisSet(key, string(record), 25*time.Hour)
}

//GetTodayLastRecord 获取上一次数据
func GetTodayLastRecord(ruleId int64) (model.StatsDays, error) {
	var result model.StatsDays
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", RuleLastPrefix, today, ruleId)
	jsonStr, err := RedisGetString(key)
	if err != nil {
		return result, err
	}
	if jsonStr == "" {
		return result, nil
	}
	return result, json.Unmarshal([]byte(jsonStr), &result)
}

func GetExistedPvKey(ruleId int64) bool {
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("%s_%s_%d", PVKeyPrefix, today, ruleId)
	existed, _ := RedisExistedKey(key)
	if existed > 0 {
		return true
	}
	return false
}

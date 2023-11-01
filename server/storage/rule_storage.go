package storage

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"redirect/model"
	"redirect/utils"
	"reflect"
	"time"
)

type Trans struct {
	Tx *sqlx.Tx
}

func GetRuleList(search, startTime, endTime string, page, size int) ([]model.Rule, error) {
	var (
		found  []model.Rule
		offset = (page - 1) * size
		query  = `SELECT  * FROM public.rules r WHERE 1=1 `
	)

	if !utils.EmptyString(search) {
		likeSearch := "%" + search + "%"
		query += fmt.Sprintf(` AND r.from_domain LIKE '%s' OR r.app_name LIKE '%s' OR r.remark LIKE '%s'`,
			likeSearch, likeSearch, likeSearch)
	}

	if !utils.EmptyString(startTime) {
		query += fmt.Sprintf(` AND r.access_time >= to_date('%s','YYYY-MM-DD')`, startTime)
	}

	if !utils.EmptyString(endTime) {
		query += fmt.Sprintf(` AND r.access_time < to_date('%s','YYYY-MM-DD')`, endTime)
	}

	query += ` ORDER BY r.id DESC LIMIT $1 OFFSET $2`
	return found, DbSelect(query, &found, size, offset)
}

func GetAllRules() ([]model.Rule, error) {
	var found []model.Rule
	query := `SELECT  * FROM public.rules r WHERE status = $1 `
	return found, DbSelect(query, &found, 1)
}

func GetRuleTotals(search, startTime, endTime string) (int, error) {
	var query = `SELECT  count(r.id) as total_count FROM public.rules r WHERE 1=1`
	if !utils.EmptyString(search) {
		likeSearch := "%" + search + "%"
		query += fmt.Sprintf(` AND r.from_domain LIKE '%s' `, likeSearch)
	}

	if !utils.EmptyString(startTime) {
		query += fmt.Sprintf(` AND r.access_time >= to_date('%s','YYYY-MM-DD')`, endTime)
	}

	if !utils.EmptyString(endTime) {
		query += fmt.Sprintf(` AND r.access_time < to_date('%s','YYYY-MM-DD')`, endTime)
	}
	var count int
	return count, DbGet(query, &count)
}

func AddRuleTrans(fromDomains []string, rules *model.Rule) error {
	tx := GetDbTx()
	var ruleId int64
	//这里有pq.Array类型,不能用LastInsertId,要用原生查询,pgsql要用RETURNING id
	query := `INSERT INTO public.rules (app_name,rule_data,status,remark,default_url,ip_blacks)` +
		` VALUES($1, $2, $3,$4,$5,$6) RETURNING id`
	err := tx.QueryRow(query, rules.AppName, rules.RuleData, rules.Status, rules.Remark, rules.DefaultUrl,
		rules.IpBlacks).Scan(&ruleId)
	if err != nil {
		tx.Rollback()
		return err
	}

	ruleDomains := getRuleDomainModel(ruleId, fromDomains)
	queryDomain := `INSERT INTO public.rule_domains (rule_id, from_domain) VALUES (:rule_id, :from_domain)`
	_, err = tx.NamedExec(queryDomain, ruleDomains)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func getRuleDomainModel(ruleId int64, fromDomains []string) []*model.RuleDomains {
	var ruleDomains []*model.RuleDomains
	for _, domain := range fromDomains {
		oneDomainRule := &model.RuleDomains{
			RuleId:     ruleId,
			FromDomain: domain,
		}
		ruleDomains = append(ruleDomains, oneDomainRule)
	}

	return ruleDomains
}

// UpdateRuleTrans 修改规则
func UpdateRuleTrans(fromDomain []string, rules *model.Rule) error {
	trans := &Trans{Tx: GetDbTx()}
	// 修改rule表数据
	err := trans.UpdateRuleTx(rules)
	if err != nil {
		trans.Tx.Rollback()
		return utils.RaiseError("修改规则失败")
	}

	//检测from_domain是否变更,未变更直接提交事务
	rawDomain := FindFromDomainsByRuleId(rules.ID)
	if true != reflect.DeepEqual(fromDomain, rawDomain) {
		err = trans.UpdateRuleDomainTx(fromDomain, rules)
		if err != nil {
			trans.Tx.Rollback()
			return err
		}
	}

	//删除缓存规则
	err = DeleteDomainCacheRules(fromDomain)
	if err != nil {
		trans.Tx.Rollback()
		return utils.RaiseError("清除缓存失败")
	}

	trans.Tx.Commit()

	return nil
}

//DelRuleTrans 删除规则
func DelRuleTrans(ids []int) error {
	var batch []*BatchQueryArgs
	//方便后边清除缓存规则
	fromDomains, err := GetAllFromDomainsByRuleIds(ids)
	if err != nil {
		return err
	}

	//删除rule
	sql1 := `DELETE FROM public.rules  WHERE  id in (?)`
	query1, args1, _ := sqlx.In(sql1, ids)
	batchQuery1 := &BatchQueryArgs{query1, args1}
	batch = append(batch, batchQuery1)

	// 删除来源域名
	sql2 := `DELETE FROM public.rule_domains  WHERE  rule_id in (?)`
	query2, args2, _ := sqlx.In(sql2, ids)
	batchQuery2 := &BatchQueryArgs{query2, args2}
	batch = append(batch, batchQuery2)

	err = DbBatchExecTx(batch)
	if err != nil {
		return err
	}

	//删除缓存
	err = DeleteDomainCacheRules(fromDomains)
	if err != nil {
		return err
	}
	return nil

}

//FindRuleById 根据规则ID查询规则
func FindRuleById(id int64) (model.Rule, error) {
	var rule model.Rule
	query := `SELECT * FROM public.rules  WHERE id = $1`
	return rule, DbGet(query, &rule, id)
}

func FindRuleWeightById(id int64) (model.RuleWeight, error) {
	var (
		rule     model.Rule
		result   model.RuleWeight
		ruleData []*model.DomainWeight
	)

	query := `SELECT * FROM public.rules  WHERE id = $1`
	err := DbGet(query, &rule, id)
	if err != nil {
		return result, err
	}

	json.Unmarshal([]byte(rule.RuleData), &ruleData)
	result = model.RuleWeight{
		ID:         rule.ID,
		AppName:    rule.AppName,
		RuleData:   ruleData,
		IpBlacks:   rule.IpBlacks,
		DefaultUrl: rule.DefaultUrl,
		Remark:     rule.Remark,
		Status:     rule.Status,
		CreatedAt:  rule.CreatedAt,
		UpdatedAt:  rule.UpdatedAt,
	}
	return result, nil
}

//GetDomainRuleFromCache 从缓存中获取规则
func GetDomainRuleFromCache(fromDomain string) (rules model.RuleWeight, err error) {
	cacheKey := utils.RuleConfig.CachePrefix + fromDomain
	expireTime := time.Duration(utils.RuleConfig.CacheExpiredTime) * time.Second
	ruleJson, err := RedisGetString(cacheKey)
	if err != nil {
		return
	}

	if !utils.EmptyString(ruleJson) {
		err = json.Unmarshal([]byte(ruleJson), &rules)
		if err != nil {
			return
		}
		return
	}

	ruleId := IsExistedFromDomain(fromDomain, 0)
	if ruleId < 1 {
		err = utils.RaiseError("rule not existed")
		return
	}

	rules, err = FindRuleWeightById(ruleId)
	if err != nil {
		return
	}

	// 缓存
	ruleBytes, _ := json.Marshal(rules)
	RedisSet(cacheKey, string(ruleBytes), expireTime)

	return
}

//DeleteDomainCacheRules 删除多个域名缓存规则
func DeleteDomainCacheRules(fromDomains []string) error {
	var cacheRuleKeys []string
	if len(fromDomains) == 0 {
		return nil
	}

	for _, domain := range fromDomains {
		fromDomain := utils.RuleConfig.CachePrefix + domain
		cacheRuleKeys = append(cacheRuleKeys, fromDomain)
	}

	//不存在的key会被忽略
	err := RedisDelete(cacheRuleKeys...)
	if err != nil {
		return err
	}
	return nil
}

//GetAllFromDomainsByRuleIds 通过规则ids获取所有域名
func GetAllFromDomainsByRuleIds(ids []int) ([]string, error) {
	var fromDomains []string
	sql := `SELECT from_domain FROM public.rule_domains  WHERE  rule_id in (?)`
	query, args, _ := sqlx.In(sql, ids)
	formatQuery := RebindQuery(query)
	return fromDomains, DbSelect(formatQuery, &fromDomains, args...)
}

//UpdateRuleTx 更新规则表事务
func (trans *Trans) UpdateRuleTx(rules *model.Rule) error {
	query := `UPDATE public.rules SET app_name = :app_name,rule_data = :rule_data,status = :status,` +
		`remark = :remark,default_url = :default_url,ip_blacks = :ip_blacks WHERE id = :id`
	_, err := trans.Tx.NamedExec(query, rules)
	if err != nil {
		return err
	}
	return nil
}

//UpdateRuleDomainTx 更新关联表rule_domains数据
func (trans *Trans) UpdateRuleDomainTx(fromDomain []string, rules *model.Rule) error {
	//先删除域名
	err := trans.deleteRuleDomainTx(rules)

	if err != nil {
		return utils.RaiseError("清空以前域名失败")
	}

	//判断是否有重复域名
	err = checkDomainExists(fromDomain, rules)
	if err != nil {
		return err
	}

	//插入域名
	err = trans.insertRuleDomainTx(fromDomain, rules)
	if err != nil {
		return utils.RaiseError("记录新的域名失败")
	}
	return nil
}

//deleteRuleDomainTx  删除域名事务
func (trans *Trans) deleteRuleDomainTx(rules *model.Rule) error {
	queryRuleDel := `DELETE FROM public.rule_domains  WHERE  rule_id = $1`
	_, err := trans.Tx.Exec(queryRuleDel, rules.ID)
	if err != nil {
		return err
	}
	return nil
}

// insertRuleDomainTx 插入域名事务
func (trans *Trans) insertRuleDomainTx(fromDomain []string, rules *model.Rule) error {
	ruleDomains := getRuleDomainModel(rules.ID, fromDomain)
	queryDomain := `INSERT INTO public.rule_domains (rule_id, from_domain) VALUES (:rule_id, :from_domain)`
	_, err := trans.Tx.NamedExec(queryDomain, ruleDomains)
	if err != nil {
		return err
	}
	return nil
}

//checkDomainExists 检测是否存在重复域名
func checkDomainExists(fromDomain []string, rules *model.Rule) error {
	for _, item := range fromDomain {
		existed := IsExistedFromDomain(item, rules.ID)
		if existed > 0 {
			return utils.RaiseError(item + " 已存在")
		}
	}
	return nil
}

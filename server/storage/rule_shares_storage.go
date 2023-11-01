package storage

import (
	"redirect/model"
)

//FindShareByRuleId 通过分享账号信息
func FindShareByRuleId(ruleId int) (model.RuleShare, error) {
	var share model.RuleShare
	query := `SELECT id,rule_id,password,share_url,created_at FROM public.rule_shares WHERE rule_id = $1`
	return share, DbGet(query, &share, ruleId)
}

//SaveRuleShareInfo 只在分享信息
func SaveRuleShareInfo(share *model.RuleShare) error {
	query := `INSERT INTO public.rule_shares (rule_id,"password",share_url) VALUES(:rule_id,:password,:share_url)`
	return DbNamedExec(query, share)
}

//UpdateRuleShareInfo 修改分享信息
func UpdateRuleShareInfo(share *model.RuleShare) error {
	query := `UPDATE public.rule_shares SET "password" = :password WHERE rule_id = :rule_id`
	return DbNamedExec(query, share)
}

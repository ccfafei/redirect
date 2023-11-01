package storage

import "redirect/model"

// IsExistedFromDomain 是否存在域名
func IsExistedFromDomain(fromDomain string, exceptId int64) int64 {
	var (
		result model.RuleDomains
		ruleId int64
	)
	query := `SELECT *  FROM public.rule_domains WHERE from_domain = $1 AND rule_id != $2`
	err := DbGet(query, &result, fromDomain, exceptId)
	if err != nil {
		return ruleId
	}
	ruleId = result.RuleId
	return ruleId
}

//FindFromDomainsByRuleId 通过规则ID查询所有域名
func FindFromDomainsByRuleId(ruleId int64) []string {
	var (
		result      []model.RuleDomains
		fromDomains []string
	)
	query := `SELECT *  FROM public.rule_domains WHERE rule_id = $1`
	err := DbSelect(query, &result, ruleId)
	if err != nil {
		return fromDomains
	}

	for _, item := range result {
		fromDomains = append(fromDomains, item.FromDomain)
	}
	return fromDomains
}

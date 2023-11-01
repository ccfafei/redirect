package service

import (
	"fmt"
	"github.com/lib/pq"
	"redirect/model"
	"redirect/storage"
	"redirect/utils"
	"redirect/utils/algorithm"
	"reflect"
)

//这里要全局初始化
var roundMap = NewRoundBalanceMap()
var weightMap = NewWeightBalanceMap()

//MatchToDomain 通过规则匹配目标域名
func MatchToDomain(fromInfo *model.FromDomainInfo, rule model.RuleWeight) string {
	var (
		toDomain   string
		defaultUrl = rule.DefaultUrl
	)

	//是否有黑名单，有的话跳到默认域名
	if len(rule.IpBlacks) > 0 {
		isIn := InIpBlacks(fromInfo.ClientIP, rule.IpBlacks)
		if isIn == true {
			return defaultUrl
		}
	}

	if len(rule.RuleData) == 0 {
		return defaultUrl
	}
	//根据权重跳转,权重默认都是100按顺序轮询跳转，否则按权重轮询
	isDefaultRule := getAllEqualWeight(rule, 100)
	if isDefaultRule {
		toDomain = RoundChoose(fromInfo.Domain, rule)
	} else {
		toDomain = WeightChoose(fromInfo.Domain, rule)
	}

	// 如果跳转域名不在时跳到默认
	if toDomain == "" {
		return defaultUrl
	}

	return toDomain
}

//GetRulesByDomain 获取域名规则
func GetRulesByDomain(fromDomain string) (rule model.RuleWeight, err error) {
	if utils.EmptyString(fromDomain) {
		err = utils.RaiseError("域名不能为空")
		return
	}

	//从缓存获取规则
	rule, err = storage.GetDomainRuleFromCache(fromDomain)
	if err != nil {
		//fmt.Println("cache error:", err)
		err = utils.RaiseError("获取域名规则失败")
		return
	}

	return
}

//InIpBlacks 是否在IP黑名单内
func InIpBlacks(ip string, ipBlacks pq.StringArray) bool {
	for _, item := range ipBlacks {
		if utils.IpRangeContains(ip, item) {
			return true
		}
	}
	return false
}

//WeightChoose 权重选择
func WeightChoose(fromDomain string, rule model.RuleWeight) string {
	ruleData := rule.RuleData
	weightData := weightMap.Get(fromDomain)

	//没有值时赋值
	if weightData == nil {
		weightData = addRawWeight(ruleData)
		weightMap.Set(fromDomain, weightData)
		weightMap.UpdatedAt = rule.UpdatedAt
		return weightData.Next()
	}

	//判断是否要更新值
	if true == checkWeightUpdated(rule, weightData) {
		weightMap.Del(fromDomain)
		rawWeight := addRawWeight(ruleData)
		weightMap.Set(fromDomain, rawWeight)
		weightMap.UpdatedAt = rule.UpdatedAt
		return rawWeight.Next()
	}

	return weightData.Next()
}

// addRawWeight 添加到权重节点中
func addRawWeight(ruleData []*model.DomainWeight) *algorithm.WeightRoundRobinBalance {
	rawWeight := &algorithm.WeightRoundRobinBalance{}
	for _, item := range ruleData {
		rawWeight.Add(item.ToDomain, fmt.Sprintf("%d", item.Weight))
	}
	return rawWeight
}

// checkWeightUpdated 判断全局权重是否更新,true要更新
func checkWeightUpdated(rule model.RuleWeight, data *algorithm.WeightRoundRobinBalance) bool {
	if rule.UpdatedAt != roundMap.UpdatedAt {
		globalRule := fetchGlobalWeightRule(data)
		if false == reflect.DeepEqual(globalRule, rule.RuleData) {
			return true
		}
	}
	return false
}

//fetchGlobalWeightRule 比较值使用
func fetchGlobalWeightRule(data *algorithm.WeightRoundRobinBalance) []model.DomainWeight {
	var globalRule []model.DomainWeight
	for _, item := range data.Rss {
		rules := model.DomainWeight{ToDomain: item.Addr, Weight: item.Weight}
		globalRule = append(globalRule, rules)
	}
	return globalRule
}

//RoundChoose 顺序轮询服务器
func RoundChoose(fromDomain string, rule model.RuleWeight) string {
	ruleData := rule.RuleData
	roundData := roundMap.Get(fromDomain)

	//没有值时
	if roundData == nil {
		roundData = addRawRound(ruleData)
		roundMap.Set(fromDomain, roundData)
		roundMap.UpdatedAt = rule.UpdatedAt
		toDomain, _ := roundData.Next()
		return toDomain
	}

	//判断是否要更新值
	if true == checkRoundUpdated(rule, roundData) {
		roundMap.Del(fromDomain)
		rawData := addRawRound(ruleData)
		roundMap.Set(fromDomain, rawData)
		roundMap.UpdatedAt = rule.UpdatedAt
		toDomain, _ := rawData.Next()
		return toDomain
	}

	toDomain, _ := roundData.Next()
	return toDomain
}

func checkRoundUpdated(rule model.RuleWeight, roundData *algorithm.Round) bool {
	if rule.UpdatedAt != roundMap.UpdatedAt {
		domains := getRuleDomains(rule.RuleData)
		if false == reflect.DeepEqual(roundData.Rss, domains) {
			return true
		}
	}
	return false
}

func getRuleDomains(ruleData []*model.DomainWeight) []string {
	var domains []string
	for _, item := range ruleData {
		domains = append(domains, item.ToDomain)
	}
	return domains
}

// getAllEqualWeight 判断默认都是不是100
func getAllEqualWeight(rule model.RuleWeight, defaultWeight int) bool {
	allEqual := true
	ruleData := rule.RuleData
	for _, item := range ruleData {
		if item.Weight != defaultWeight {
			allEqual = false
			break
		}
	}
	return allEqual
}

func addRawRound(ruleData []*model.DomainWeight) *algorithm.Round {
	var rawRound = &algorithm.Round{}
	for _, item := range ruleData {
		rawRound.Add(item.ToDomain)
	}
	return rawRound
}

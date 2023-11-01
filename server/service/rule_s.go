package service

import (
	"encoding/json"
	"fmt"
	"redirect/model"
	"redirect/storage"
	"redirect/utils"
	"strings"
)

//GetRuleList 查询规则列表
func GetRuleList(search, startTime, endTime string, page, size int) (data model.PageInfo, err error) {
	var AllRuleData []*model.ResponseRule
	if page < 1 || size < 1 {
		return
	}
	rules, err := storage.GetRuleList(search, startTime, endTime, page, size)
	if err != nil {
		fmt.Println("rule list err:", err)
		err = utils.RaiseError("内部错误，请联系管理员")
		return
	}

	for _, item := range rules {
		var jsonData []*model.DomainWeight
		var fromDomains []string
		json.Unmarshal([]byte(item.RuleData), &jsonData)
		fromDomains = storage.FindFromDomainsByRuleId(item.ID)
		ruleOne := &model.ResponseRule{
			ID:         item.ID,
			AppName:    item.AppName,
			FromDomain: fromDomains,
			RuleData:   jsonData,
			Status:     item.Status,
			Remark:     item.Remark,
			IpBlacks:   item.IpBlacks,
			DefaultUrl: item.DefaultUrl,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		}
		AllRuleData = append(AllRuleData, ruleOne)
	}

	total, err := storage.GetRuleTotals(search, startTime, endTime)
	if err != nil {
		err = utils.RaiseError("内部错误，请联系管理员")
		return
	}

	data = model.PageInfo{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  AllRuleData,
	}

	return
}

// AddRule 添加规则
func AddRule(param model.AddRuleParam) error {
	var (
		ruleData string
	)
	//判断来源域名是否存在
	fromDomains := param.FromDomain
	err := checkAddFromDomain(fromDomains, 0)
	if err != nil {
		return err
	}

	if len(param.RuleData) > 0 {
		jsonData, _ := json.Marshal(param.RuleData)
		ruleData = string(jsonData)
	}

	rules := &model.Rule{
		AppName:    param.AppName,
		RuleData:   ruleData,
		DefaultUrl: param.DefaultUrl,
		IpBlacks:   param.IpBlacks,
		Remark:     param.Remark,
		Status:     param.Status,
	}

	err = storage.AddRuleTrans(fromDomains, rules)
	if err != nil {
		fmt.Println("add rule error:", err)
		return utils.RaiseError("添加规则失败")
	}

	return nil
}

//UpdateRule 修改规则
func UpdateRule(param *model.UpdateRuleParam) error {
	found, err := storage.FindRuleById(param.ID)
	if err != nil {
		return utils.RaiseError("查找数据失败")
	}
	if found.ID == 0 {
		return utils.RaiseError("未查找到数据")
	}

	if !utils.EmptyString(param.AppName) {
		found.AppName = param.AppName
	}

	fromDomains := param.FromDomain

	var ruleData string
	if len(param.RuleData) > 0 {
		jsonData, _ := json.Marshal(param.RuleData)
		ruleData = string(jsonData)
		found.RuleData = ruleData
	}

	if param.Status != -1 {
		found.Status = param.Status
	}

	if !utils.EmptyString(param.Remark) {
		found.Remark = param.Remark
	}

	if param.DefaultUrl != "" {
		found.DefaultUrl = param.DefaultUrl
	}

	if len(param.IpBlacks) > 0 {
		found.IpBlacks = param.IpBlacks
	}

	err = storage.UpdateRuleTrans(fromDomains, &found)
	if err != nil {
		return utils.RaiseError(err.Error())
	}

	return nil
}

//DeleteRule 删除规则
func DeleteRule(strIds string) error {
	intIds := utils.SplitStrIdsToInt(strIds, ",")
	if len(intIds) == 0 {
		return utils.RaiseError("ids格式不对")
	}

	err := storage.DelRuleTrans(intIds)
	if err != nil {
		fmt.Println("delete error:", err)
		return utils.RaiseError("删除失败")
	}

	return nil
}

//GetRuleInfo 获取规则信息
func GetRuleInfo(ruleId int64) (model.Rule, error) {
	return storage.FindRuleById(ruleId)
}

// checkAddFromDomain 判断域名是否存在
func checkAddFromDomain(fromDomains []string, exceptRuleId int64) error {
	for _, item := range fromDomains {
		fromDomain := strings.TrimSpace(item)
		ruleId := storage.IsExistedFromDomain(fromDomain, 0)
		if ruleId > 0 {
			return utils.RaiseError(fromDomain + "已存在")
		}
	}
	return nil
}

package storage

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetDomainRuleFromCache(t *testing.T) {
	init4Test(t)
	fromDomain := "dnsce.com"
	t.Run("TestGetDomainRuleFromCache", func(t *testing.T) {
		rules, err := GetDomainRuleFromCache(fromDomain)
		if err != nil {
			t.Error(err)
			return
		}
		if rules.ID == 0 {
			t.Error("获取规则失败")
			return
		}

		if rules.RuleData != "" {
			var ruleData interface{}
			err = json.Unmarshal([]byte(rules.RuleData), &ruleData)
			if err != nil {
				t.Error("解析json失败")
				return
			}
			fmt.Println("ruleData:", ruleData)
		}

	})
}

func TestGetAllFromDomainsByRuleIds(t *testing.T) {
	init4Test(t)
	t.Run("TestGetFromDomain", func(t *testing.T) {
		ids := []int{24}
		domains, err := GetAllFromDomainsByRuleIds(ids)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(domains)
	})
}

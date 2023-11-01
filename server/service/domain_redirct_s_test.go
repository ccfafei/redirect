package service

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"redirect/model"
	"redirect/storage"
	"testing"
)

func TestMatchToDomain(t *testing.T) {
	init4Test(t)
	t.Run("MatchToDomain", func(t *testing.T) {
		// 构造数据
		userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36"
		referer := ""
		uvCookie := "uv_" + faker.UUIDDigit()
		fromInfo := &model.FromDomainInfo{
			Domain:    "t1.com",
			ClientIP:  "10.1.1.9",
			UserAgent: userAgent,
			Referer:   referer,
			UvCookie:  uvCookie,
		}
		//获取规则
		rule, err := GetRulesByDomain(fromInfo.Domain)
		if err != nil {
			fmt.Println("获取规则失败")
			return
		}

		//匹配域名
		toDomain := MatchToDomain(fromInfo, rule)
		fmt.Println("跳转的地址为:", toDomain)
		return
	})
}

func TestCacheDomain(t *testing.T) {
	init4Test(t)
	t.Run("MatchToDomain", func(t *testing.T) {
		fromDomain := "t1.com"
		rule, err := storage.GetDomainRuleFromCache(fromDomain)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(rule)
	})
}

func TestInIpBlacks(t *testing.T) {
	ipBlacks := []string{
		"111.175.80.210",
		"36.56.0.0/24",
		"39.8.0.1-39.8.0.254",
	}
	t.Run("TestInIpBlacks", func(t *testing.T) {
		ip1 := "111.175.80.210"
		rs1 := InIpBlacks(ip1, ipBlacks)
		fmt.Printf("%s是否在IP黑名单中?%v \n", ip1, rs1)

		ip2 := "36.56.10.10"
		rs2 := InIpBlacks(ip2, ipBlacks)
		fmt.Printf("%s是否在IP黑名单中?%v \n", ip2, rs2)

		ip3 := "36.56.0.10"
		rs3 := InIpBlacks(ip3, ipBlacks)
		fmt.Printf("%s是否在IP黑名单中?%v \n", ip3, rs3)

		ip4 := "39.8.0.18"
		rs4 := InIpBlacks(ip4, ipBlacks)
		fmt.Printf("%s是否在IP黑名单中?%v \n", ip4, rs4)
	})
}

func TestRoundChoose(t *testing.T) {
	var ruleData []model.DomainWeight
	init4Test(t)
	t.Run("TestRoundChoose", func(t *testing.T) {
		fromDomain := "t1.com"
		rule, err := storage.GetDomainRuleFromCache(fromDomain)
		if err != nil {
			t.Error(err)
		}
		json.Unmarshal([]byte(rule.RuleData), &ruleData)
		toDomain := RoundChoose(fromDomain, ruleData)
		fmt.Println(toDomain)
	})
}

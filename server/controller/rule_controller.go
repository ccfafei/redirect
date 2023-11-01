package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"net/url"
	"redirect/core"
	"redirect/model"
	"redirect/service"
	"redirect/utils"
	"strconv"
	"strings"
)

//GetRuleList 根据类型查询所有规则
func GetRuleList(ctx *gin.Context) {
	strPage := ctx.DefaultQuery("page", strconv.Itoa(DefaultPageNum))
	strSize := ctx.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	search := ctx.DefaultQuery("search", "")
	startTime := ctx.DefaultQuery("start_time", "")
	endTime := ctx.DefaultQuery("end_time", "")
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DefaultPageNum
	}

	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DefaultPageSize
	}

	rules, err := service.GetRuleList(search, startTime, endTime, page, size)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest("查询失败"))
		return
	}
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(rules))
}

//AddRules 添加规则
func AddRules(ctx *gin.Context) {
	var addRuleParam model.AddRuleParam
	//validate 验证
	err := ctx.ShouldBindJSON(&addRuleParam)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}

	// 来源域名格式
	fromDomain, err := validateFromDomain(addRuleParam.FromDomain)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}
	addRuleParam.FromDomain = fromDomain

	//验证IP黑名单
	err = validateIpBlacks(addRuleParam.IpBlacks)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}

	// 验证跳转域名
	rule, err := validateToDomain(addRuleParam.RuleData)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}
	addRuleParam.RuleData = rule

	//验证默认域名
	defaultUrl, err := validateDefaultURL(addRuleParam.DefaultUrl)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}
	addRuleParam.DefaultUrl = defaultUrl

	err = service.AddRule(addRuleParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest(fmt.Sprintf("%v", err)))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

//UpdateRule 修改规则
func UpdateRule(ctx *gin.Context) {
	var updateRuleParam *model.UpdateRuleParam
	//validate 验证
	err := ctx.ShouldBindJSON(&updateRuleParam)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}

	//验证IP黑名单
	err = validateIpBlacks(updateRuleParam.IpBlacks)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}

	// 验证跳转域名
	rule, err := validateToDomain(updateRuleParam.RuleData)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}
	updateRuleParam.RuleData = rule

	//验证默认域名
	defaultUrl, err := validateDefaultURL(updateRuleParam.DefaultUrl)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}
	updateRuleParam.DefaultUrl = defaultUrl
	fromDomain, err := validateFromDomain(updateRuleParam.FromDomain)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}
	updateRuleParam.FromDomain = fromDomain

	err = service.UpdateRule(updateRuleParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest(fmt.Sprintf("%v", err)))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

//DeleteRule 删除规则
func DeleteRule(ctx *gin.Context) {
	var deleteParam model.DeleteRuleParam
	err := ctx.ShouldBindJSON(&deleteParam)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}

	err = service.DeleteRule(deleteParam.Ids)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

func GetRuleInfo(ctx *gin.Context) {
	ruleIdStr := ctx.Query("rule_id")
	if ruleIdStr == "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("rule_id不能为空"))
		return
	}

	ruleId, _ := strconv.Atoi(ruleIdStr)
	ruleInfo, err := service.GetRuleInfo(int64(ruleId))
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("获取数据失败"))
		return
	}
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(ruleInfo))
}

//validateIpBlacks 验证黑白名单
func validateIpBlacks(ipArr pq.StringArray) error {
	if len(ipArr) > 0 {
		for _, ip := range ipArr {
			if false == utils.ValidateIP(ip) {
				return utils.RaiseError(ip + "格式不正确")
			}
		}
	}
	return nil
}

func validateFromDomain(fromDomains []string) (result []string, err error) {
	if len(fromDomains) > 0 {
		for _, item := range fromDomains {
			fromDomain := utils.TrimString(item)
			if fromDomain == "" {
				err = utils.RaiseError("来源域名不能为空")
				return
			}
			//给fromDomain加前缀
			if !strings.HasPrefix(fromDomain, "http://") && !strings.HasPrefix(fromDomain, "https://") {
				fromDomain = "http://" + fromDomain
			}
			u, err1 := url.Parse(fromDomain)
			if err1 != nil {
				err = utils.RaiseError("域名格式不对")
				return
			}
			result = append(result, u.Host)
		}
	}
	return
}

func validateToDomain(ruleData []*model.DomainWeight) (result []*model.DomainWeight, err error) {
	if len(ruleData) > 0 {
		for _, item := range ruleData {
			toDomain := utils.TrimString(item.ToDomain)
			if toDomain == "" {
				continue
			}
			if !strings.HasPrefix(toDomain, "http://") && !strings.HasPrefix(toDomain, "https://") {
				err = utils.RaiseError("目标网址必须以http://或https://开头")
				return
			}
			_, err = url.Parse(toDomain)
			if err != nil {
				err = utils.RaiseError("网址格式不正确")
				return
			}
			item.ToDomain = toDomain
			result = append(result, item)
		}
	}
	return
}

func validateDefaultURL(defaultUrl string) (result string, err error) {
	result = utils.TrimString(defaultUrl)
	if result != "" {
		_, err = url.Parse(result)
		if err != nil {
			err = utils.RaiseError("默认网址不正确")
			return
		}
	}
	return
}

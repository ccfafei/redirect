package service

import (
	"fmt"
	"redirect/model"
	"redirect/storage"
	"redirect/utils"
	"strings"
	"time"
)

// ShareLogin 登录
func ShareLogin(ruleId int, password string) (model.RuleShare, error) {
	var found model.RuleShare
	found, err := storage.FindShareByRuleId(ruleId)
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员")
	}

	if found.ID < 1 {
		return found, utils.RaiseError("未查找到数据")
	}

	res := utils.AesEncrypt(password, utils.ShareConfig.DesKey)
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员")
	}

	if !strings.EqualFold(found.Password, res) {
		return found, utils.RaiseError("用户名或密码错误")
	}

	return found, nil
}

// NewShareToken 生成分享token
func NewShareToken(share model.RuleShare) (model.ShareLoginResult, error) {
	var result model.ShareLoginResult
	account := fmt.Sprintf("guest_%d", share.ID)
	token, err := utils.GenerateToken(account, share.RuleId)
	if err != nil {
		return result, err
	}

	jwtExpireTime := time.Duration(utils.JwtConfig.JwtExpiredTime) * time.Second
	jwtKey := utils.JwtConfig.ShareJwtPrefix + account
	err = storage.RedisSet(jwtKey, token, jwtExpireTime)
	if err != nil {
		return result, err
	}
	result = model.ShareLoginResult{RuleId: share.RuleId, Token: token}
	return result, nil
}

//GetShareInfo 获取分享信息
func GetShareInfo(ruleId int) (model.RuleShare, error) {
	var (
		result model.RuleShare
		err    error
	)
	result, err = storage.FindShareByRuleId(ruleId)
	if err != nil {
		return result, err
	}
	//为空创建默认密码
	if result.IsEmpty() {
		defaultPwd := utils.RandSimplePassword()
		result, err = createRuleShare(ruleId, defaultPwd)
		if err != nil {
			return result, err
		}
	}
	if result.Password != "" {
		result.Password = utils.AesDecrypt(result.Password, utils.ShareConfig.DesKey)
	}
	if result.ShareUrl != "" {
		result.ShareUrl = fmt.Sprintf("%s%s", utils.ShareConfig.ShareDomain, result.ShareUrl)
	}

	return result, err
}

//UpdateShareRuleInfo 修改分享密码
func UpdateShareRuleInfo(param model.ShareParam) error {
	password := utils.AesEncrypt(param.Password, utils.ShareConfig.DesKey)
	share := &model.RuleShare{RuleId: param.RuleId, Password: password}
	found, err := storage.FindShareByRuleId(param.RuleId)
	if err != nil {
		return utils.RaiseError("内部错误，请联系管理员")
	}
	if found.IsEmpty() {
		return utils.RaiseError("未查询到内容")
	}

	err = storage.UpdateRuleShareInfo(share)
	if err != nil {
		return utils.RaiseError("修改失败，请稍后再试")
	}

	return nil
}

//createRuleShare 创建默认分享信息,并显示
func createRuleShare(ruleId int, defaultPwd string) (model.RuleShare, error) {
	var result model.RuleShare
	shareUrl := fmt.Sprintf("/#/login?rule_id=%d", ruleId)
	password := utils.AesEncrypt(defaultPwd, utils.ShareConfig.DesKey)
	share := &model.RuleShare{RuleId: ruleId, Password: password, ShareUrl: shareUrl}
	err := storage.SaveRuleShareInfo(share)
	if err != nil {
		return result, err
	}
	result, err = storage.FindShareByRuleId(ruleId)
	if err != nil {
		return result, err
	}
	return result, nil
}

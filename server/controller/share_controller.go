package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redirect/core"
	"redirect/model"
	"redirect/service"
	"redirect/utils"
	"strconv"
)

//DoShareLogin 分享登录
func DoShareLogin(ctx *gin.Context) {
	var shareLoginParam model.ShareParam
	//validate 验证
	ctx.ShouldBindJSON(&shareLoginParam)
	if shareLoginParam.RuleId < 1 {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("任务链接地址不完整，请核对信息!"))
		return
	}

	if shareLoginParam.Password == "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("密码不能为空!"))
		return
	}

	// 用户名密码有效性验证
	loginShare, err := service.ShareLogin(shareLoginParam.RuleId, shareLoginParam.Password)
	if err != nil || loginShare.IsEmpty() {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("密码错误,请重新输入"))
		return
	}

	//生成token,保存在redis
	tokenData, err := service.NewShareToken(loginShare)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("系统错误，请与客服联系"))
		return
	}
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(tokenData))
}

//GetShareUrl 获取分享信息
func GetShareUrl(ctx *gin.Context) {
	strRuleId := ctx.DefaultQuery("rule_id", "")
	ruleId, _ := strconv.Atoi(strRuleId)
	share, err := service.GetShareInfo(ruleId)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("获取分享信息失败，请联系客服"))
		return
	}
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(share))
}

//UpdateShareUrl 修改分享信息
func UpdateShareUrl(ctx *gin.Context) {
	var param model.ShareParam
	//validate 验证
	err := ctx.ShouldBindJSON(&param)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}
	err = service.UpdateShareRuleInfo(param)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

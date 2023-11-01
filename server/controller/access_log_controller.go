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

//GetAccessLogList 根据类型查询所有规则
func GetAccessLogList(ctx *gin.Context) {
	strPage := ctx.DefaultQuery("page", strconv.Itoa(DefaultPageNum))
	strSize := ctx.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	strRuleId := ctx.DefaultQuery("rule_id", "")
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

	ruleId, _ := strconv.Atoi(strRuleId)
	accessLogs, err := service.GetPagedAccessLogs(int64(ruleId), search, startTime, endTime, page, size)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("查询失败"))
		return
	}
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(accessLogs))
}

//DeleteAccessLogs 删除日志
func DeleteAccessLogs(ctx *gin.Context) {
	var deleteLogParam model.DeleteLogParam
	err := ctx.ShouldBindJSON(&deleteLogParam)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}

	err = service.DeleteAccessLogByIds(deleteLogParam.Ids)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

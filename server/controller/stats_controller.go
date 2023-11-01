package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redirect/service"
	"redirect/utils"
	"strconv"
)

//GetAllDataNum 统计所有总数
func GetAllDataNum(ctx *gin.Context) {
	ruleIdStr := ctx.Query("rule_id")
	if ruleIdStr == "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("任务ID不能为空,请联系客服"))
		return
	}
	ruleId, _ := strconv.Atoi(ruleIdStr)
	result := service.TotalAllData(int64(ruleId))
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(result))
}

//GetChartData 统计图表数据
func GetChartData(ctx *gin.Context) {
	ruleIdStr := ctx.Query("rule_id")
	if ruleIdStr == "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("任务ID不能为空,请联系客服"))
		return
	}
	dateType := ctx.Query("date_type")
	if dateType == "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("日期类型不能为空"))
		return
	}
	ruleId, _ := strconv.Atoi(ruleIdStr)
	result := service.TotalChartData(int64(ruleId), dateType)
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(result))
}

//GetTopRank 统计当天排名数据
func GetTopRank(ctx *gin.Context) {
	ruleIdStr := ctx.Query("rule_id")
	if ruleIdStr == "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("任务ID不能为空,请联系客服"))
		return
	}
	topStr := ctx.DefaultQuery("top", "10")

	ruleId, _ := strconv.Atoi(ruleIdStr)
	top, _ := strconv.Atoi(topStr)
	result := service.GetTopRank(int64(ruleId), top)
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(result))
}

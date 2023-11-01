package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"redirect/controller"
	"redirect/core"
	"redirect/utils"
)

func InitRouteAdmin() (http.Handler, error) {

	if utils.AppConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	//跨域
	router.Use(core.Cors())

	//日志
	logFilePath := utils.LogConfig.LogFilePath
	logFileName := "admin.log"
	fileName := path.Join(logFilePath, logFileName)
	router.Use(core.LoggerToFile(fileName))

	//账号登录
	router.POST("/login", controller.DoLogin)
	router.GET("/captcha", controller.RequestCaptchaImage)
	router.GET("/captcha/:imageId", controller.ServeCaptchaImage) //暂时不用，为方便看验证码

	//后台验证的api
	api := router.Group("/api", core.AdminAuth())
	//用户管理
	api.GET("/account/all", controller.GetAllAdmin)
	api.GET("/account/info", controller.GetAdminInfo)
	api.POST("/account/add", controller.AddAdmin)
	api.POST("/account/update", controller.UpdateAdmin)
	api.POST("/account/del", controller.DelAdmin)
	api.POST("/logout", controller.DoLogout)

	// 规则管理
	api.GET("/rule/all", controller.GetRuleList)
	api.POST("/rule/add", controller.AddRules)
	api.POST("/rule/update", controller.UpdateRule)
	api.POST("/rule/delete", controller.DeleteRule)

	//日志管理
	api.GET("/logs/all", controller.GetAccessLogList)
	api.POST("/logs/delete", controller.DeleteAccessLogs)

	//分享管理
	api.GET("/share/info", controller.GetShareUrl)
	api.POST("/share/update", controller.UpdateShareUrl)

	// #################################################################################
	//不鉴权
	router.POST("/share_login", controller.DoShareLogin)
	//数据统计，单独鉴权
	share := router.Group("/user", core.ShareAuth())
	share.GET("/rule/info", controller.GetRuleInfo)
	share.GET("/stats/total", controller.GetAllDataNum)
	share.GET("/stats/chart", controller.GetChartData)
	share.GET("/stats/rank", controller.GetTopRank)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, utils.ResultJsonError("404"))
		return
	})
	return router, nil
} // end of router01

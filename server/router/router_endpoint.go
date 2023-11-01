package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redirect/controller"
	"redirect/core"
	"redirect/utils"
)

// InitRouteEndPoint 跳转程序
func InitRouteEndPoint() (http.Handler, error) {

	if utils.AppConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	//跨域
	router.Use(core.Cors())

	//阻止缓存
	router.Use(core.NoCache())

	//路由
	router.GET("/", controller.DoDomainRedirect)

	//统计
	router.Static("/assets", "./assets")
	router.GET("/stat", controller.DoWebStats)

	router.NoRoute(func(ctx *gin.Context) {
		return
	})

	return router, nil
} // end of router

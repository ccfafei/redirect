package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"redirect/storage"
	"redirect/utils"
	"strings"
	"time"
)

const (
	authorizationHeaderKey = "token"
	userId                 = "user_id"
)

//Cors 跨域处理
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

//AdminAuth 后端认证
func AdminAuth() gin.HandlerFunc {
	adminJwtKey := utils.JwtConfig.AdminJwtPrefix
	return jwtAuth(adminJwtKey)
}

//ShareAuth 分享认证
func ShareAuth() gin.HandlerFunc {
	adminJwtKey := utils.JwtConfig.ShareJwtPrefix
	return jwtAuth(adminJwtKey)
}

//jwtAuth 基于JWT的认证中间件
func jwtAuth(redisPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中取出
		signToken := c.Request.Header.Get(authorizationHeaderKey)
		if signToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResultJsonUnauthorized("Token不能为空"))
			return
		}
		// 校验token
		myClaims, err := utils.ParseToken(signToken)
		if err != nil {
			fmt.Println("parse error:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResultJsonUnauthorized("认证失败，请重新登录"))
			return
		}
		//过期时间
		if time.Now().Unix() > myClaims.ExpiresAt {
			c.AbortWithStatusJSON(http.StatusRequestTimeout, utils.ResultJsonUnauthorized("登录超时，请重新登录"))
			return
		}
		// 查询缓存是否存在
		jwtKey := redisPrefix + myClaims.Account
		cacheToken, err := storage.RedisGetString(jwtKey)
		if err != nil || cacheToken == "" {
			c.AbortWithStatusJSON(http.StatusRequestTimeout, utils.ResultJsonUnauthorized("Token is not Exited"))
			return
		}

		c.Set(userId, myClaims.Id)
		c.Next()
	}
}

//NoCache 阻止缓存响应
func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		ctx.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		ctx.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		ctx.Next()
	}
}

// WebLogFormatHandler Customized log format for web
func WebLogFormatHandler(server string) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if !strings.HasPrefix(param.Path, "/assets") {
			return fmt.Sprintf("[%s | %s] %s %s %d %s \t%s %s %s \n",
				server,
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		} // end of if
		return ""
	}) // end of formatter
} // end of func

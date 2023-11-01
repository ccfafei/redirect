package controller

import (
	"bytes"
	"encoding/base64"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"redirect/core"
	"redirect/model"
	"redirect/service"
	"redirect/utils"
)

const (
	DefaultPageNum  = 1
	DefaultPageSize = 20
)

// DoLogin 登录
func DoLogin(ctx *gin.Context) {
	var loginParam model.LoginParam

	//validate 验证
	err := ctx.ShouldBindJSON(&loginParam)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}

	// 验证码有效性验证
	if utils.CaptchaConfig.Enable == true {
		if !captcha.VerifyString(loginParam.CaptchaId, loginParam.CaptchaText) {
			ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("验证码错误，请刷新页面再重新尝试！"))
			return
		}
	}

	// 用户名密码有效性验证
	loginUser, err := service.Login(loginParam.Account, loginParam.Password)
	if err != nil || loginUser.IsEmpty() {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("密码错误"))
		return
	}

	//生成token,并缓存
	tokenData, err := service.NewToken(loginUser)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest("系统错误，生成token失败"))
		return
	}
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(tokenData))
}

// ServeCaptchaImage 生成验证码
func ServeCaptchaImage(c *gin.Context) {
	captcha.Server(utils.CaptchaConfig.Width, utils.CaptchaConfig.Height).ServeHTTP(c.Writer, c.Request)
}

// RequestCaptchaImage 获取验证码:验证码id和base64图片
func RequestCaptchaImage(ctx *gin.Context) {
	captchaId := captcha.New()
	var image bytes.Buffer
	err := captcha.WriteImage(&image, captchaId, utils.CaptchaConfig.Width, utils.CaptchaConfig.Height)
	if err != nil {
		return
	}
	imageData := base64.StdEncoding.EncodeToString([]byte(image.String()))
	data := map[string]string{
		"captcha_id":   captchaId,
		"captcha_data": "data:image/png;base64," + imageData,
	}
	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(data))
}

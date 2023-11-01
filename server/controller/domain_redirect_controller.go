package controller

import (
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"redirect/model"
	"redirect/service"
	"strings"
	"time"
)

//DoDomainRedirect 获取域名跳转
func DoDomainRedirect(ctx *gin.Context) {
	//识别域名信息
	fromInfo := GetDomainInfo(ctx)
	fromDomain := fromInfo.Domain
	if fromDomain == "" {
		//fmt.Println("获取来源域名失败")
		NotFound(ctx)
		return
	}

	//获取规则
	rule, err := service.GetRulesByDomain(fromInfo.Domain)
	if err != nil {
		NotFound(ctx)
		return
	}

	//匹配域名
	toDomain := service.MatchToDomain(fromInfo, rule)
	if toDomain == "" {
		NotFound(ctx)
		return
	}

	//统计pv uv ip
	err = service.NewPvUvIpData(fromInfo, rule.ID, toDomain)
	if err != nil {
		NotFound(ctx)
		return
	}

	// 写入日志
	err = service.NewAccessLog(fromInfo, rule.ID, toDomain)
	if err != nil {
		NotFound(ctx)
		return
	}

	//fmt.Printf("域名%s正在跳到到%s...\n", fromDomain, toDomain)

	//开始跳转
	ctx.Redirect(http.StatusMovedPermanently, toDomain)
	return

}

//DoWebStats 网站分析统计
func DoWebStats(ctx *gin.Context) {
	//识别域名信息
	fromInfo := GetRequestInfo(ctx)
	fromDomain := fromInfo.Domain
	if fromDomain == "" {
		//fmt.Println("获取来源域名失败")
		NotFound(ctx)
		return
	}
	fmt.Println("网站统计:", fromInfo)

	//获取规则
	rule, err := service.GetRulesByDomain(fromInfo.Domain)
	if err != nil {
		//fmt.Println("获取规则失败")
		NotFound(ctx)
		return
	}

	//统计pv uv ip
	err = service.NewPvUvIpData(fromInfo, rule.ID, "")
	if err != nil {
		//fmt.Println("统计pv uv ip 失败")
		NotFound(ctx)
		return
	}

	// 写入日志
	err = service.NewAccessLog(fromInfo, rule.ID, "")
	if err != nil {
		//fmt.Println("写入日志失败")
		NotFound(ctx)
		return
	}
}

// NotFound 404定义
func NotFound(ctx *gin.Context) {
	ctx.AbortWithError(200, errors.New("404 not found"))
}

//GetDomainInfo 获取来源域名信息
func GetDomainInfo(ctx *gin.Context) *model.FromDomainInfo {
	var info *model.FromDomainInfo
	uvCookie, _ := GetUvCookie(ctx)
	hosts := strings.Split(ctx.Request.Host, ":")
	if len(hosts) == 0 {
		return info
	}
	info = &model.FromDomainInfo{
		Domain:    hosts[0],
		ClientIP:  ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
		Referer:   ctx.Request.Referer(),
		UvCookie:  uvCookie,
	}
	return info
}

func GetRequestInfo(ctx *gin.Context) *model.FromDomainInfo {
	var info *model.FromDomainInfo
	domain := ctx.DefaultQuery("domain", "")
	userAgent := ctx.DefaultQuery("agent", "")
	referrer := ctx.DefaultQuery("referrer", "")
	uvCookie := ctx.DefaultQuery("cookie", "")

	info = &model.FromDomainInfo{
		Domain:    domain,
		ClientIP:  ctx.ClientIP(),
		UserAgent: userAgent,
		Referer:   referrer,
		UvCookie:  uvCookie,
	}
	return info
}

// GetUvCookie 生成uv cookie
func GetUvCookie(ctx *gin.Context) (string, error) {
	const UvCookie = "uv_cookie"
	cookies, err := ctx.Cookie(UvCookie)
	maxAge := getMaxAge()
	if err != nil {
		cookieUUID := "uv_" + faker.UUIDDigit()
		ctx.SetCookie(UvCookie, cookieUUID, maxAge, "/", "", false, true)
		cookies, err = ctx.Cookie(UvCookie)
		return cookies, err
	}
	return cookies, nil
}

func getMaxAge() int {
	currentTime := time.Now().Unix()
	lastTimeStr := fmt.Sprintf("%s 23:59:59", time.Now().Format("2006-01-02"))
	t, _ := time.Parse("2006-01-02 15:04:05", lastTimeStr)
	lastTime := t.Unix()
	maxAge := lastTime - currentTime
	if maxAge == 0 {
		maxAge = 86400
	}
	return int(maxAge)
}

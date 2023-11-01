package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"redirect/core"
	"redirect/model"
	"redirect/service"
	"redirect/utils"
	"strconv"
	"strings"
)

//GetAdminInfo 查询某个用户信息
func GetAdminInfo(ctx *gin.Context) {
	userId := ctx.Query("id")
	if utils.EmptyString(userId) {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest("用户ID不能为空"))
		return
	}

	uid, _ := strconv.Atoi(strings.TrimSpace(userId))
	info, err := service.GetUserInfo(uid)
	info.Password = "******"
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest("查询失败"))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(info))
}

//GetAllAdmin 获取管理员列表
func GetAllAdmin(ctx *gin.Context) {
	strPage := ctx.DefaultQuery("page", strconv.Itoa(DefaultPageNum))
	strSize := ctx.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	search := ctx.DefaultQuery("search", "")
	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = DefaultPageNum
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		size = DefaultPageSize
	}

	data, err := service.FindAllUsers(search, page, size)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest("获取用户列表失败"))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccessWithData(data))
}

//AddAdmin 添加管理员
func AddAdmin(ctx *gin.Context) {
	var addAdminParam model.AddAdminParam
	//validate 验证
	err := ctx.ShouldBindJSON(&addAdminParam)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}
	user := &model.User{Account: addAdminParam.Account, Password: addAdminParam.Password, Name: addAdminParam.Name}
	err = service.NewUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest(fmt.Sprintf("%v", err)))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

// UpdateAdmin 修改管理员
func UpdateAdmin(ctx *gin.Context) {
	var updateAdminParam model.UpdateAdminParam
	//validate 验证
	err := ctx.ShouldBindJSON(&updateAdminParam)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}

	//uid, _ := strconv.Atoi(strings.TrimSpace(updateAdminParam.Id))
	user := &model.User{ID: updateAdminParam.ID, Password: updateAdminParam.Password, Name: updateAdminParam.Name}
	err = service.UpdateAdmin(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest("修改失败"))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

// DelAdmin 删除管理
func DelAdmin(ctx *gin.Context) {
	var deleteAdminParam model.DeleteAdminParam
	err := ctx.ShouldBindJSON(&deleteAdminParam)
	validateErr := core.GetFirstValidateError(err)
	if validateErr != "" {
		ctx.JSON(http.StatusOK, utils.ResultJsonBadRequest(validateErr))
		return
	}

	err = service.DelUserByIds(deleteAdminParam.Ids)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResultJsonBadRequest(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

// DoLogout 登出
func DoLogout(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusOK, utils.ResultJsonError("用户已退出"))
		return
	}

	uid, _ := strconv.Atoi(fmt.Sprintf("%v", userId))
	err := service.Logout(uid)
	if err != nil {
		c.JSON(http.StatusOK, utils.ResultJsonError("系统错误"))
		return
	}
	c.JSON(http.StatusOK, utils.ResultJsonSuccess())
}

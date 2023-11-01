package service

import (
	"fmt"
	"redirect/model"
	"strings"
	"time"

	"redirect/storage"
	"redirect/utils"
)

const AdminUserPrefix = "redirectAdmin#"
const AdminCookiePrefix = "redirectCookie#"

// Login 登录
func Login(account string, password string) (model.User, error) {
	var found model.User
	found, err := storage.FindUserByAccount(account)
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员")
	}

	if found.IsEmpty() {
		return found, utils.RaiseError("用户名或密码错误")
	}

	res, err := utils.PasswordBase58Hash(password)
	if err != nil {
		return found, utils.RaiseError("内部错误，请联系管理员")
	}

	if !strings.EqualFold(found.Password, res) {
		return found, utils.RaiseError("用户名或密码错误")
	}

	return found, nil
}

//FindAllUsers 获取所有用户列表
func FindAllUsers(search string, page, size int) (data model.PageInfo, err error) {

	if page < 1 || size < 1 {
		return
	}
	users, err := storage.FindAllUsers(search, page, size)
	if err != nil {
		err = utils.RaiseError("内部错误，请联系管理员")
		return
	}
	total, err := storage.FindAllUsersTotals(search)
	if err != nil {
		err = utils.RaiseError("内部错误，请联系管理员")
		return
	}

	data = model.PageInfo{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  users,
	}

	return
}

//UpdateAdmin 修改管理员信息(密码，名称等)
func UpdateAdmin(user *model.User) error {
	found, err := storage.FindUserByUserId(user.ID)
	if err != nil {
		return err
	}

	if found.IsEmpty() {
		return utils.RaiseError("未查询到账号")
	}

	// 空密码时不破坏原始密码
	if !utils.EmptyString(user.Password) {
		np, err := utils.PasswordBase58Hash(user.Password)
		if err != nil {
			return err
		}

		found.Password = np
	}

	if !utils.EmptyString(user.Name) {
		found.Name = user.Name
	}

	err = storage.UpdateUser(found)
	if err != nil {
		return err
	}

	return nil
}

//NewUser 创建账号
func NewUser(user *model.User) error {
	found, err := storage.FindUserByAccount(user.Account)
	if err != nil {
		return utils.RaiseError("系统错误，查询失败")
	}

	if !found.IsEmpty() {
		return utils.RaiseError(fmt.Sprintf("用户名 %s 已存在", user.Account))
	}

	// 密码加密
	user.Password, _ = utils.PasswordBase58Hash(user.Password)
	err = storage.NewUser(user)
	if err != nil {
		return utils.RaiseError("系统错误，创建失败")
	}

	return nil
}

//GetUserInfo 获取用户信息
func GetUserInfo(id int) (found model.User, err error) {
	found, err = storage.FindUserByUserId(id)
	if err != nil {
		return
	}

	if found.IsEmpty() {
		err = utils.RaiseError("该用户不存在")
		return
	}

	return
}

// DelUserByIds 通过ids删除用户
func DelUserByIds(strIds string) error {
	intIds := utils.SplitStrIdsToInt(strIds, ",")
	if len(intIds) == 0 {
		return utils.RaiseError("ids格式不对")
	}

	user, _ := storage.FindUserByAccount("admin")
	if utils.InArray(user.ID, intIds) {
		return utils.RaiseError("admin不能删除")
	}

	err := storage.DelUsers(intIds)
	if err != nil {
		//fmt.Println("delete error:", err)
		return utils.RaiseError("删除失败")
	}

	return nil
}

// NewToken 生成token
func NewToken(user model.User) (result *model.LoginResult, err error) {
	token, err := utils.GenerateToken(user.Account, user.ID)
	if err != nil {
		return
	}
	result = &model.LoginResult{
		Id:      user.ID,
		Account: user.Account,
		Token:   token,
	}

	jwtExpireTime := time.Duration(utils.JwtConfig.JwtExpiredTime) * time.Second
	jwtKey := utils.JwtConfig.AdminJwtPrefix + user.Account
	err = storage.RedisSet(jwtKey, token, jwtExpireTime)
	if err != nil {
		return
	}

	return
}

//Logout 退出，得新生成token
func Logout(id int) (err error) {
	found, err := GetUserInfo(id)
	if err != nil {
		return err
	}
	//设置token为空并过期
	jwtKey := utils.JwtConfig.AdminJwtPrefix + found.Account
	storage.RedisSet(jwtKey, "", time.Second)
	return
}

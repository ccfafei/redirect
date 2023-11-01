package model

import (
	"reflect"
)

// User 用户
type User struct {
	ID       int    `db:"id" json:"id"`
	Account  string `db:"account" json:"account"`
	Password string `db:"password" json:"password"`
	Name     string `db:"name" json:"name"`
}

type LoginResult struct {
	Id      int    `json:"id"`
	Account string `json:"account"`
	Token   string `json:"token"`
}

// IsEmpty 判断是否为空
func (user User) IsEmpty() bool {
	return reflect.DeepEqual(user, User{})
}

type LoginParam struct {
	Account     string `json:"account" binding:"required"`
	Password    string `json:"password" binding:"required,min=6,max=18"`
	CaptchaText string `json:"captcha_text" binding:"required"`
	CaptchaId   string `json:"captcha_id" binding:"required"`
}

type AddAdminParam struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=18"`
	Name     string `json:"name"`
}

type UpdateAdminParam struct {
	ID       int    `json:"id" binding:"required"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type DeleteAdminParam struct {
	Ids string `json:"ids" binding:"required"`
}

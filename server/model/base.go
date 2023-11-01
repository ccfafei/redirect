package model

import (
	"github.com/golang-jwt/jwt"
	"time"
)

//Claims jwt 使用
type Claims struct {
	Account string `json:"account"`
	Id      int    `json:"id"`
	jwt.StandardClaims
}

//PageInfo 分页
type PageInfo struct {
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Data  interface{} `json:"data"`
}

// ResultJson 返回结果
type ResultJson struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
	Date    time.Time   `json:"date"`
}

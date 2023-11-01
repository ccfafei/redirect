package model

import (
	"reflect"
	"time"
)

type RuleShare struct {
	ID        int       `db:"id" json:"id"`
	RuleId    int       `db:"rule_id" json:"rule_id"`
	Password  string    `db:"password" json:"password"`
	ShareUrl  string    `db:"share_url" json:"share_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ShareParam struct {
	RuleId   int    `db:"rule_id" json:"rule_id" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=18"`
}

type ShareLoginResult struct {
	RuleId int    `db:"rule_id" json:"rule_id" binding:"required"`
	Token  string `db:"token" json:"token"`
}

// IsEmpty 判断是否为空
func (share RuleShare) IsEmpty() bool {
	return reflect.DeepEqual(share, RuleShare{})
}

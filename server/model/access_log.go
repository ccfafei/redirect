package model

import (
	"time"
)

// AccessLog 访问日志
type AccessLog struct {
	ID         int64     `db:"id" json:"id"`
	LogUUID    string    `db:"log_uuid" json:"log_uuid"`
	RuleID     int64     `db:"rule_id" json:"rule_id"`
	FromDomain string    `db:"from_domain" json:"from_domain"`
	ToDomain   string    `db:"to_domain" json:"to_domain"`
	AccessTime time.Time `db:"access_time" json:"access_time"`
	Ip         *string   `db:"ip" json:"ip,omitempty"`
	UserAgent  *string   `db:"user_agent" json:"user_agent,omitempty"`
	Referer    *string   `db:"referer" json:"referer,omitempty"`
	UvCookie   *string   `db:"uv_cookie" json:"uv_cookie,omitempty"`
}

type DeleteLogParam struct {
	Ids string `json:"ids" binding:"required"`
}

// FromDomainInfo 域名来源信息
type FromDomainInfo struct {
	Domain    string
	ClientIP  string
	UserAgent string
	Referer   string
	UvCookie  string
}

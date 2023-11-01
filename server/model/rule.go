package model

import (
	"github.com/lib/pq"
	"time"
)

//Rule 规则
type Rule struct {
	ID         int64          `db:"id" json:"id"`
	AppName    string         `db:"app_name" json:"app_name"`
	RuleData   string         `db:"rule_data" json:"rule_data"`
	DefaultUrl string         `db:"default_url" json:"default_url"`
	IpBlacks   pq.StringArray `db:"ip_blacks" json:"ip_blacks"`
	Remark     string         `db:"remark" json:"remark"`
	Status     int            `db:"status" json:"status"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updated_at"`
}

type RuleWeight struct {
	ID         int64           `db:"id" json:"id"`
	AppName    string          `db:"app_name" json:"app_name"`
	RuleData   []*DomainWeight `db:"rule_data" json:"rule_data"`
	DefaultUrl string          `db:"default_url" json:"default_url"`
	IpBlacks   pq.StringArray  `db:"ip_blacks" json:"ip_blacks"`
	Remark     string          `db:"remark" json:"remark"`
	Status     int             `db:"status" json:"status"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
}

// ResponseRule 返回json
type ResponseRule struct {
	ID         int64           `db:"id" json:"id"`
	AppName    string          `db:"app_name" json:"app_name"`
	FromDomain []string        `db:"from_domain" json:"from_domain"`
	RuleData   []*DomainWeight `db:"rule_data" json:"rule_data"`
	DefaultUrl string          `db:"default_url" json:"default_url"`
	IpBlacks   pq.StringArray  `db:"ip_blacks" json:"ip_blacks"`
	Status     int             `db:"status" json:"status"`
	Remark     string          `db:"remark" json:"remark"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
}

type DomainWeight struct {
	ToDomain string `json:"to_domain"`
	Weight   int    `json:"weight"`
}

type AddRuleParam struct {
	AppName    string          `json:"app_name" binding:"required"`
	FromDomain []string        `json:"from_domain" binding:"required"`
	RuleData   []*DomainWeight `json:"rule_data"`
	DefaultUrl string          `json:"default_url"`
	IpBlacks   pq.StringArray  `json:"ip_blacks"`
	Status     int             `json:"status"`
	Remark     string          `json:"remark"`
}

type DeleteRuleParam struct {
	Ids string `json:"ids" binding:"required"`
}

type UpdateRuleParam struct {
	ID         int64           `json:"id" binding:"required"`
	AppName    string          `json:"app_name"`
	FromDomain []string        `json:"from_domain"`
	RuleData   []*DomainWeight `json:"rule_data"`
	DefaultUrl string          `json:"default_url"`
	IpBlacks   pq.StringArray  `json:"ip_blacks"`
	Status     int             `json:"status"`
	Remark     string          `json:"remark"`
}

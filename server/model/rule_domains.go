package model

import "time"

type RuleDomains struct {
	ID         int64     `db:"id" json:"id"`
	RuleId     int64     `db:"rule_id" json:"rule_id"`
	FromDomain string    `db:"from_domain" json:"from_domain"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

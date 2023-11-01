package model

import "time"

type StatsMinutes struct {
	ID         int64     `db:"id" json:"id"`
	RuleId     int64     `db:"rule_id" json:"rule_id"`
	AccessTime time.Time `db:"access_time" json:"access_time"`
	PvNum      int64     `db:"pv_num" json:"pv_num"`
	UvNum      int64     `db:"uv_num" json:"uv_num"`
	IpNum      int64     `db:"ip_num" json:"ip_num"`
}

type TopRank struct {
	FromDomain string `db:"from_domain" json:"from_domain"`
	Num        int    `db:"num" json:"num"`
}

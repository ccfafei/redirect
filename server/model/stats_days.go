package model

import "time"

type StatsDays struct {
	ID         int       `db:"id" json:"id"`
	RuleId     int       `db:"rule_id" json:"rule_id"`
	AccessTime time.Time `db:"access_time" json:"access_time"`
	PvNum      int64     `db:"pv_num" json:"pv_num"`
	UvNum      int64     `db:"uv_num" json:"uv_num"`
	IpNum      int64     `db:"ip_num" json:"ip_num"`
}

type TotalStats struct {
	AccessTime string `db:"access_time" json:"access_time"`
	PvNum      int64  `db:"pv_num" json:"pv_num"`
	UvNum      int64  `db:"uv_num" json:"uv_num"`
	IpNum      int64  `db:"ip_num" json:"ip_num"`
}

type TotalNum struct {
	TodayPvNum     int64 `json:"today_pv_num"`
	TodayUvNum     int64 `json:"today_uv_num"`
	TodayIpNum     int64 `json:"today_ip_num"`
	YesterdayPvNum int64 `json:"yesterday_pv_num"`
	YesterdayUvNum int64 `json:"yesterday_uv_num"`
	YesterdayIpNum int64 `json:"yesterday_ip_num"`
}

type TotalChartData struct {
	PvData []*BaseData `json:"pv_data"`
	UvData []*BaseData `json:"uv_data"`
	IpData []*BaseData `json:"ip_data"`
}

type BaseData struct {
	Ts  string `json:"ts"`
	Num int64  `json:"num"`
}

package core

import (
	"github.com/robfig/cron/v3"
	"redirect/service"
	"time"
)

func RunJob() {

	//nyc, _ := time.LoadLocation("Asia/Shanghai")
	cron := cron.New(cron.WithSeconds()) // 设置定时任务时区

	//同步redis日志到数据库,并清理相关数据
	cron.AddFunc("*/10 * * * * ?", func() {
		service.StoreAccessLogs()
	})

	//定时统计5分钟日志数据
	cron.AddFunc("0 */5 * * * ?", func() {
		//fmt.Println("正在统计分钟数据")
		service.JobStatsMinutes()
	})

	//定时统计昨天数据,晚上12点 "0 0 0 * * ?"
	cron.AddFunc("0 0 0 * * ?", func() {
		service.JobStatsDays()
	})

	// 零点10分清理日志
	cron.AddFunc("0 10 0 * * ?", func() {
		service.DeleteAccessLogHistory(-3)
	})

	cron.Start()
	defer cron.Stop()

	t := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t.C:
			t.Reset(time.Second * 10)
		}
	}
}

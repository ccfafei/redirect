package service

import (
	"errors"
	"redirect/model"
	"redirect/storage"
	"time"
)

//JobStatsMinutes 定时统计分钟数据
func JobStatsMinutes() error {
	rules, err := storage.GetAllRules()
	if err != nil {
		return err
	}

	stats, err := getStatsMinutes(rules)
	if err != nil {
		return err
	}

	err = storage.InsertStatsMinutes(stats)
	if err != nil {
		return err
	}

	//删除3天前的数据
	err = storage.DeleteStatsMinutes()
	if err != nil {
		return err
	}

	return nil
}

// getStatsMinutes
func getStatsMinutes(rules []model.Rule) ([]model.StatsMinutes, error) {
	var stats []model.StatsMinutes
	for _, item := range rules {
		isExisted := storage.GetExistedPvKey(item.ID)
		if isExisted == false {
			continue
		}
		oneStats, err1 := diffInsertData(item.ID)
		//fmt.Println("one stats:", oneStats)
		if err1 != nil {
			continue
		}
		if oneStats.PvNum == 0 && oneStats.UvNum == 0 && oneStats.IpNum == 0 {
			continue
		}
		stats = append(stats, oneStats)
	}
	return stats, nil
}

//diffInsertData 计算增长的数据
func diffInsertData(ruleId int64) (model.StatsMinutes, error) {
	var result model.StatsMinutes

	//先获取上次历史数据
	lastData, err := storage.GetTodayLastRecord(ruleId)
	if err != nil {
		//fmt.Println("last data:", err)
		return result, err
	}

	today := time.Now().Format("2006-01-02")
	totalData, err := storage.TotalTodayNum(ruleId, today)
	if err != nil {
		return result, err
	}
	if totalData.PvNum == 0 && totalData.UvNum == 0 && totalData.IpNum == 0 {
		return result, errors.New("no data")
	}

	diffPvNum := totalData.PvNum - lastData.PvNum
	diffUvNum := totalData.UvNum - lastData.UvNum
	diffIpNum := totalData.IpNum - lastData.IpNum
	if diffUvNum < 0 || diffUvNum < 0 || diffIpNum < 0 {
		return result, errors.New("error format")
	}
	oneStats := model.StatsMinutes{
		RuleId:     ruleId,
		AccessTime: time.Now(),
		PvNum:      diffPvNum,
		UvNum:      diffUvNum,
		IpNum:      diffIpNum,
	}

	//最后更新上一次数据
	storage.SaveTodayLastRecord(ruleId, totalData)
	return oneStats, nil
}

//JobStatsDays 定时统计前一天数据,删除前3个月
func JobStatsDays() error {
	rules, err := storage.GetAllRules()
	if err != nil {
		return err
	}
	for _, item := range rules {
		//存在先跳过
		fund, _ := storage.GetYesterdaysByRuleId(item.ID)
		if fund.ID > 0 {
			continue
		}

		//查询数据
		result, err1 := getYeterDayData(item.ID)
		if err1 != nil {
			continue
		}

		//插入表
		err = storage.InsertStatsDays(result)
		if err != nil {
			return err
		}

	}

	//删除3月前的数据
	err = storage.DeleteStatsDays()
	if err != nil {
		return err
	}
	return nil
}

func getYeterDayData(ruleId int64) (model.StatsDays, error) {
	var stats model.StatsDays
	var yesterday = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	accessTime, _ := time.ParseInLocation("2006-01-02", yesterday, time.Local)
	//从缓存中查询
	result, err := storage.TotalTodayNum(ruleId, yesterday)
	if err != nil {
		return stats, err
	}

	stats = model.StatsDays{
		RuleId:     int(ruleId),
		AccessTime: accessTime,
		PvNum:      result.PvNum,
		UvNum:      result.UvNum,
		IpNum:      result.IpNum,
	}

	return stats, nil
}

//TotalAllData 统计总数
func TotalAllData(ruleId int64) *model.TotalNum {
	today := time.Now().Format("2006-01-02")
	result, err := storage.TotalTodayNum(ruleId, today)
	if err != nil {
		return nil
	}
	yesterday, err := storage.GetYesterdaysByRuleId(ruleId)
	if err != nil {
		return nil
	}
	totalAllNums := &model.TotalNum{
		TodayUvNum:     result.UvNum,
		TodayPvNum:     result.PvNum,
		TodayIpNum:     result.IpNum,
		YesterdayUvNum: yesterday.PvNum,
		YesterdayPvNum: yesterday.UvNum,
		YesterdayIpNum: yesterday.IpNum,
	}
	return totalAllNums
}

//TotalChartData 统计图表数据
func TotalChartData(ruleId int64, dateType string) model.TotalChartData {
	var data model.TotalChartData
	switch dateType {
	case "today":
		beforeDay := 0
		data = getChartData(beforeDay, ruleId, "hours")
	case "yesterday":
		beforeDay := -1
		data = getChartData(beforeDay, ruleId, "hours")
	case "week":
		beforeDay := -7
		data = getChartData(beforeDay, ruleId, "days")
	case "month":
		beforeDay := -30
		data = getChartData(beforeDay, ruleId, "days")
	}

	return data
}

//GetTopRank 获取排名靠前数据
func GetTopRank(ruleId int64, top int) []model.TopRank {
	var (
		data []model.TopRank
		err  error
	)
	data, err = storage.GetTopRank(ruleId, top)
	if err != nil {
		return data
	}
	return data
}

//getChartData  获取图表数据,比较关键是补数据
func getChartData(beforeDay int, ruleId int64, types string) model.TotalChartData {
	// 获取数据
	chartData := getTypesChartData(beforeDay, ruleId, types)

	//组合数组数据
	newData := getMapChartByTimeKey(chartData, types)

	//将日期起始转数组
	dateArr := getRangeDate(beforeDay, types)

	//补齐数据
	data := padChartData(dateArr, newData)

	return data
}

//padChartData 补齐数据
func padChartData(dateArr []int64, newData map[int64]model.TotalStats) model.TotalChartData {
	var pvDataArr, uvDataArr, ipDataArr []*model.BaseData
	for _, h := range dateArr {
		accessTime := time.Unix(h, 0).Format("2006-01-02 15:04:05")
		var defaultBase = &model.BaseData{Ts: accessTime, Num: 0}
		pvData, uvData, ipData := defaultBase, defaultBase, defaultBase
		_, ok := newData[h]
		if ok {
			item := newData[h]
			pvData = &model.BaseData{Ts: accessTime, Num: item.PvNum}
			uvData = &model.BaseData{Ts: accessTime, Num: item.UvNum}
			ipData = &model.BaseData{Ts: accessTime, Num: item.IpNum}

		}
		pvDataArr = append(pvDataArr, pvData)
		uvDataArr = append(uvDataArr, uvData)
		ipDataArr = append(ipDataArr, ipData)
	}

	data := model.TotalChartData{PvData: pvDataArr, UvData: uvDataArr, IpData: ipDataArr}
	return data
}

//getMapChartByTimeKey 以时间为key组合数组
func getMapChartByTimeKey(chartData []model.TotalStats, types string) map[int64]model.TotalStats {
	loc, _ := time.LoadLocation("Local")
	var newData = make(map[int64]model.TotalStats)
	for _, item := range chartData {
		var stamp time.Time
		if types == "hours" {
			stamp, _ = time.ParseInLocation("2006-01-02 15", item.AccessTime, loc)
		} else {
			stamp, _ = time.ParseInLocation("2006-01-02", item.AccessTime, loc)
		}
		key := stamp.Unix()
		newData[key] = item
	}
	return newData
}

//getTypesChartData 根据类型获取相关数据
func getTypesChartData(beforeDay int, ruleId int64, types string) []model.TotalStats {
	var (
		chartData []model.TotalStats
		err       error
	)
	date := time.Now().AddDate(0, 0, beforeDay).Format("2006-01-02")
	if types == "hours" {
		chartData, err = storage.GetHoursChart(date, ruleId)
		if err != nil {
			return chartData
		}
	} else if types == "days" {
		chartData, err = storage.GetDaysChart(date, ruleId)
		if err != nil {
			return chartData
		}
	} else {
		return chartData
	}
	return chartData
}

// getRangeDate 生成时间数组
func getRangeDate(beforeDay int, types string) []int64 {
	var dateArr []int64
	if types == "hours" {
		dateArr = getHourArr(beforeDay)
	} else if types == "days" {
		dateArr = getDayArr(beforeDay)
	} else {
	}
	return dateArr
}

func getDayArr(beforeDay int) []int64 {
	var arr []int64
	t := time.Now()
	nowTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	for i := beforeDay; i < 0; i++ {
		dy := nowTime + int64(i*86400)
		arr = append(arr, dy)
	}
	return arr
}

func getHourArr(beforeDay int) []int64 {
	var arr []int64
	local, _ := time.LoadLocation("Local")
	date := time.Now().AddDate(0, 0, beforeDay).Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", date, local)
	nowTime := t.Unix()
	h := 24
	if beforeDay == 0 {
		h = time.Now().Hour()
	}
	for i := 0; i <= h; i++ {
		dh := nowTime + int64(i*3600)
		arr = append(arr, dh)
	}
	return arr
}

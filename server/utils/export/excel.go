package export

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"redirect/model"
	"strconv"
)

func AccessLogToExcel(logs []model.AccessLog) ([]byte, error) {
	if logs == nil {
		return nil, errors.New("数据为空")
	}
	f := excelize.NewFile()
	index, _ := f.NewSheet("Sheet1")
	// 填充表头
	f.SetCellValue("Sheet1", "A1", "短链接")
	f.SetCellValue("Sheet1", "B1", "访问时间")
	f.SetCellValue("Sheet1", "C1", "访问IP")
	f.SetCellValue("Sheet1", "D1", "UserAgent")
	for i := 0; i < len(logs); i++ {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), logs[i].FromDomain)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), logs[i].AccessTime)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), logs[i].Ip.String)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), logs[i].UserAgent.String)
	}
	f.SetActiveSheet(index)
	if excellBytes, erorrW := f.WriteToBuffer(); erorrW != nil {
		return nil, erorrW
	} else {
		return excellBytes.Bytes(), nil
	}

}

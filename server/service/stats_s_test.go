package service

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJobStatsMinutes(t *testing.T) {
	init4Test(t)
	if err := JobStatsMinutes(); err != nil {
		t.Error(err)
	}
}

func TestTotalChartData(t *testing.T) {
	init4Test(t)
	data := TotalChartData(24, "today")
	rs, _ := json.Marshal(data)
	fmt.Println(string(rs))
}

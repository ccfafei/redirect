package storage

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"math/rand"
	"redirect/model"
	"testing"
	"time"
)

func randomTime() time.Time {
	min := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func TestInsertAccessLogs(t *testing.T) {
	init4Test(t)
	var logs []model.AccessLog
	for i := 0; i < 10; i++ {
		ip := faker.IPv4()
		userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36"
		referer := ""
		uvCookie := "uv_" + faker.UUIDDigit()
		log := model.AccessLog{
			RuleID:     1,
			FromDomain: faker.DomainName(),
			ToDomain:   faker.DomainName(),
			Ip:         &ip,
			AccessTime: randomTime(),
			UserAgent:  &userAgent,
			Referer:    &referer,
			UvCookie:   &uvCookie,
		}
		logs = append(logs, log)
	}
	t.Run("TestInsertAccessLogs", func(t *testing.T) {
		err := InsertAccessLogs(logs)
		fmt.Println("log error:", err)
		if err != nil {
			t.Error(err)
			return
		}
	})
}

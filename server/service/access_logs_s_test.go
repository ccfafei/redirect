package service

import (
	"redirect/model"
	"testing"

	"github.com/bxcodec/faker/v3"

	"redirect/storage"
	"redirect/utils"
)

func TestStoreAccessLog(t *testing.T) {
	init4Test(t)
	if err := StoreAccessLogs(); err != nil {
		t.Error(err)
	}
}

func TestNewAccessLog(t *testing.T) {
	init4Test(t)
	for i := 0; i < 10; i++ {
		fromInfo := &model.FromDomainInfo{
			Domain:    faker.DomainName(),
			ClientIP:  faker.IPv4(),
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36",
			Referer:   "",
			UvCookie:  "",
		}
		if err := NewAccessLog(fromInfo, 1, "baidu.com"); err != nil {
			t.Error(err)
		}
	}
}

func init4Test(t *testing.T) {
	_, err := utils.InitConfig("../config.ini")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = storage.InitDatabaseService()
	if err != nil {
		t.Error(err)
		return
	}

	_, err = storage.InitRedisService()
	if err != nil {
		t.Error(err)
		return
	}
}

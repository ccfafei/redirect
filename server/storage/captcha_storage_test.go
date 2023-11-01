//
// Custom redis storage for captcha, Accroding to https://github.com/dchest/captcha/blob/master/store.go
//
// An object implementing Store interface can be registered with SetCustomStore
// function to handle storage and retrieval of captcha ids and solutions for
// them, replacing the default memory store.
//
// It is the responsibility of an object to delete expired and used captchas
// when necessary (for example, the default memory store collects them in Set
// method after the certain amount of captchas has been stored.)

package storage

import (
	"fmt"
	"testing"
	"time"

	"redirect/utils"

	"github.com/dchest/captcha"
)

func TestNewRedisStore(t *testing.T) {

	_, err := utils.InitConfig("../config.ini")
	if err != nil {
		t.Error(err)
		return
	}

	rs, err := InitRedisService()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("TestNewRedisStore", func(t *testing.T) {
		s := CaptchaRedisStore{KeyPrefix: "capTest", Expiration: 2 * time.Minute, RedisService: rs}

		captcha.SetCustomStore(&s)

		id := captcha.New()
		fmt.Println(id)

	})

}

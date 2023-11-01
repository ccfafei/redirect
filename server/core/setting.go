package core

import (
	"fmt"
	"github.com/dchest/captcha"
	"redirect/storage"
	"redirect/utils"
	"strings"
	"time"
)

func InitSettings(cmdConfig string) error {
	// Things MUST BE DONE before app starts
	_, err := utils.InitConfig(cmdConfig)
	if err != nil {
		fmt.Println("Config initialization failed.", err)
		return err
	}

	rs, err := storage.InitRedisService()
	if err != nil {
		fmt.Println("Redis initialization failed.", err)
		return err
	}

	if strings.EqualFold("redis", strings.ToLower(utils.CaptchaConfig.Store)) {
		expiration := time.Duration(utils.CaptchaConfig.CacheExpiredTime) * time.Second
		crs := storage.CaptchaRedisStore{KeyPrefix: utils.CaptchaConfig.CachePrefix, Expiration: expiration, RedisService: rs}
		captcha.SetCustomStore(&crs)
	}

	_, err = storage.InitDatabaseService()
	if err != nil {
		fmt.Println("database initialization failed.", err)
		return err
	}
	return nil
}

package service

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"log"
	"redirect/model"
	"redirect/storage"
	"redirect/utils"
	"time"
)

const accessLogsPrefix = "REDIRECT_ACCESS_LOGS#"
const blackLogsKey = "black_insert_keys"
const batchRedisSize = 5000
const batchSqlSize = 1000

// NewAccessLog 记录访问日志
func NewAccessLog(fromInfo *model.FromDomainInfo, ruleId int64, toDomain string) error {
	logUUID := faker.UUIDDigit()
	accessLog := &model.AccessLog{
		RuleID:     ruleId,
		LogUUID:    logUUID,
		FromDomain: fromInfo.Domain,
		ToDomain:   toDomain,
		AccessTime: time.Now(),
		Ip:         &fromInfo.ClientIP,
		UserAgent:  &fromInfo.UserAgent,
		Referer:    &fromInfo.Referer,
		UvCookie:   &fromInfo.UvCookie,
	}

	logJson, _ := json.Marshal(accessLog)
	key := fmt.Sprintf("%s%s", accessLogsPrefix, logUUID)
	err := storage.RedisSet30m(key, logJson)
	if err != nil {
		log.Println(err)
		return utils.RaiseError("内部错误，请联系管理员")
	}

	return nil
}

// StoreAccessLogs 将访问日志存入数据库
func StoreAccessLogs() error {
	keys, err := storage.RedisScan4Keys(accessLogsPrefix + "*")
	if err != nil {
		return err
	}
	batchKeys := keys
	if len(keys) > batchRedisSize {
		batchKeys = keys[0:batchRedisSize]
	}

	//重组数据
	logs, err := createLogsModelFromRedisKey(batchKeys)
	if err != nil {
		return err
	}

	// 插入数据库
	err = storage.BatchInsertAccessLogs(logs, batchSqlSize)
	if err != nil {
		return err
	}

	// 删除keys
	err = storage.RedisDelete(batchKeys...)
	if err != nil {
		return err
	}

	return nil
}

// createLogsModelFromRedisKey 将keys中数据组合
func createLogsModelFromRedisKey(keys []string) ([]model.AccessLog, error) {
	var logs []model.AccessLog
	if len(keys) == 0 {
		err := utils.RaiseError("no key")
		return logs, err
	}

	for _, k := range keys {
		v, err := storage.RedisGetString(k)
		if err != nil {
			continue
		}
		accessLog := model.AccessLog{}
		json.Unmarshal([]byte(v), &accessLog)
		logs = append(logs, accessLog)
	}

	return logs, nil
}

// filterBlackListKey 过滤黑名单的keys,返回有效值
func filterBlackListKey(insertKeys []string) ([]string, error) {
	var diffKeys []string
	if len(insertKeys) == 0 {
		err := utils.RaiseError("empty keys")
		return insertKeys, err
	}

	blackKeys := storage.RedisSMembers(blackLogsKey)
	if len(blackKeys) == 0 {
		diffKeys = insertKeys //差集为新插入的值
	} else {
		diffKeys = utils.DifferenceStringsArr(insertKeys, blackKeys) //比较求差集
	}

	if len(diffKeys) == 0 {
		err := utils.RaiseError("no add keys")
		return diffKeys, err
	}

	//加入到黑名单
	err := storage.RedisSAdd(blackLogsKey, diffKeys...)
	if err != nil {
		return diffKeys, err
	}

	return diffKeys, nil

}

// GetPagedAccessLogs 获取分页访问日志
func GetPagedAccessLogs(ruleId int64, search string, start, end string, page, size int) (data model.PageInfo, err error) {
	if page < 1 || size < 1 {
		return
	}

	allAccessLogs, err := storage.FindAllAccessLogs(ruleId, search, start, end, page, size)
	if err != nil {
		err = utils.RaiseError("内部错误，请联系管理员")
		return
	}

	total, err := storage.FindAllAccessPageTotal(ruleId, search, start, end)
	if err != nil {
		err = utils.RaiseError("内部错误，请联系管理员")
		return
	}

	data = model.PageInfo{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  allAccessLogs,
	}

	return
}

//DeleteAccessLogByIds 删除日志
func DeleteAccessLogByIds(strIds string) error {
	intIds := utils.SplitStrIdsToInt(strIds, ",")
	if len(intIds) == 0 {
		return utils.RaiseError("ids格式不对")
	}

	err := storage.DeleteAccessLogs(intIds)
	if err != nil {
		return utils.RaiseError("删除失败")
	}

	return nil
}

//DeleteAccessLogHistory 删除历史数据
func DeleteAccessLogHistory(day int) error {
	beforeTime := time.Now().AddDate(0, 0, day).Format("2006-01-02")
	err := storage.DeleteAccessLogsHistory(beforeTime)
	if err != nil {
		return err
	}
	return nil
}

//NewPvUvIpData 生成pv,uv,ip
func NewPvUvIpData(fromInfo *model.FromDomainInfo, ruleId int64, toDomain string) error {

	err := storage.SaveRulesPV(ruleId)
	if err != nil {
		return err
	}

	err = storage.SaveRulesUV(ruleId, fromInfo.Domain, fromInfo.UvCookie)
	if err != nil {
		return err
	}

	err = storage.SaveRulesIP(ruleId, fromInfo.Domain, fromInfo.ClientIP)
	if err != nil {
		return err
	}
	return nil
}

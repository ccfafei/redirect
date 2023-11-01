package utils

import (
	"strings"

	"gopkg.in/ini.v1"
)

const Version = "1.9"

var (
	DatabaseConfig     DatabaseConfigInfo
	AppConfig          AppConfigInfo
	RedisConfig        RedisConfigInfo
	RedisClusterConfig RedisClusterConfigInfo
	CaptchaConfig      CaptchaConfigInfo
	JwtConfig          JwtConfigInfo
	RuleConfig         RuleConfigInfo
	ShareConfig        ShareConfigInfo
	LogConfig          LogConfigInfo
)

type CaptchaConfigInfo struct {
	Enable           bool
	Store            string
	Width            int
	Height           int
	CachePrefix      string
	CacheExpiredTime int
}

type JwtConfigInfo struct {
	JwtKey         string
	JwtExpiredTime int
	AdminJwtPrefix string
	ShareJwtPrefix string
}

// AppConfigInfo 应用配置
type AppConfigInfo struct {
	EndpointPort    int
	AdminPort       int
	UrlPrefix       string
	Debug           bool
	WebReadTimeout  int
	WebWriteTimeout int
}

type LogConfigInfo struct {
	LogFilePath string
}

// RedisClusterConfigInfo redis配置
type RedisClusterConfigInfo struct {
	Hosts    []string
	User     string
	Password string
	PoolSize int
}

// RedisConfigInfo redis配置
type RedisConfigInfo struct {
	Host     string
	User     string
	Password string
	Database int
	PoolSize int
}

// DatabaseConfigInfo 数据库配置
type DatabaseConfigInfo struct {
	Host         string
	Port         int
	User         string
	Password     string
	DbName       string
	MaxOpenConns int
	MaxIdleConn  int
}

type RuleConfigInfo struct {
	CachePrefix      string
	CacheExpiredTime int
}

type ShareConfigInfo struct {
	ShareDomain string
	DesKey      string
}

// InitConfig 初始化配置
func InitConfig(file string) (*ini.File, error) {
	cfg, err := ini.Load(file)
	if err != nil {
		return nil, nil
	}

	section := cfg.Section("postgres")
	DatabaseConfig.Host = section.Key("host").String()
	DatabaseConfig.Port = section.Key("port").MustInt()
	DatabaseConfig.MaxOpenConns = section.Key("max_open_conn").MustInt()
	DatabaseConfig.MaxIdleConn = section.Key("max_idle_conn").MustInt()
	DatabaseConfig.User = section.Key("user").String()
	DatabaseConfig.Password = section.Key("password").String()
	DatabaseConfig.DbName = section.Key("database").String()

	appSection := cfg.Section("app")
	AppConfig.Debug = appSection.Key("debug").MustBool()
	AppConfig.EndpointPort = appSection.Key("endpoint_port").MustInt()
	AppConfig.AdminPort = appSection.Key("admin_port").MustInt()
	AppConfig.UrlPrefix = appSection.Key("url_prefix").String()
	AppConfig.WebReadTimeout = appSection.Key("web_read_timeout").MustInt()
	AppConfig.WebWriteTimeout = appSection.Key("web_write_timeout").MustInt()

	redisSection := cfg.Section("redis")
	RedisConfig.Host = redisSection.Key("host").String()
	RedisConfig.User = redisSection.Key("username").String()
	RedisConfig.Password = redisSection.Key("password").String()
	RedisConfig.Database = redisSection.Key("database").MustInt()
	RedisConfig.PoolSize = redisSection.Key("pool_size").MustInt()

	redisClusterSection := cfg.Section("redis-cluster")
	hosts := redisClusterSection.Key("hosts").String()
	if !EmptyString(hosts) {
		hostsArr := strings.Split(hosts, ",")
		RedisClusterConfig.Hosts = hostsArr
	}
	RedisClusterConfig.User = redisClusterSection.Key("username").String()
	RedisClusterConfig.Password = redisClusterSection.Key("password").String()
	RedisClusterConfig.PoolSize = redisClusterSection.Key("pool_size").MustInt()

	captchaSection := cfg.Section("captcha")
	CaptchaConfig.Enable = captchaSection.Key("enable").MustBool()
	CaptchaConfig.Store = captchaSection.Key("store").String()
	CaptchaConfig.Width = captchaSection.Key("width").MustInt()
	CaptchaConfig.Height = captchaSection.Key("height").MustInt()
	CaptchaConfig.CachePrefix = captchaSection.Key("captcha_cache_prefix").String()
	CaptchaConfig.CacheExpiredTime = captchaSection.Key("captcha_expired_time").MustInt()

	jwtSection := cfg.Section("jwt")
	JwtConfig.JwtKey = jwtSection.Key("jwt_key").String()
	JwtConfig.JwtExpiredTime = jwtSection.Key("jwt_expired_time").MustInt()
	JwtConfig.AdminJwtPrefix = jwtSection.Key("admin_jwt_prefix").String()
	JwtConfig.ShareJwtPrefix = jwtSection.Key("share_jwt_prefix").String()

	ruleSection := cfg.Section("rules")
	RuleConfig.CachePrefix = ruleSection.Key("rule_cache_prefix").String()
	RuleConfig.CacheExpiredTime = ruleSection.Key("rule_expired_time").MustInt()

	shareSection := cfg.Section("share")
	ShareConfig.ShareDomain = shareSection.Key("share_domain").String()
	ShareConfig.DesKey = shareSection.Key("des_key").String()

	LogSection := cfg.Section("log")
	LogConfig.LogFilePath = LogSection.Key("log_file_path").String()

	return cfg, err
}

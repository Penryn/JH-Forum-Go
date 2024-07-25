// 该代码文件定义了应用程序的配置参数和设置初始化过程。

package conf

import (
	"log"
	"time"

	"github.com/alimy/tryst/cfg"
)

// loggerSetting 定义了日志配置
var loggerSetting *loggerConf

// loggerFileSetting 定义了文件日志配置
var loggerFileSetting *loggerFileConf

// loggerMeiliSetting 定义了 Meili 日志配置
var loggerMeiliSetting *loggerMeiliConf


// redisSetting 定义了 Redis 配置
var redisSetting *redisConf

// PyroscopeSetting 定义了 Pyroscope 配置
var PyroscopeSetting *pyroscopeConf

// DatabaseSetting 定义了数据库配置
var DatabaseSetting *databaseConf

// MysqlSetting 定义了 MySQL 数据库配置
var MysqlSetting *mysqlConf

// PprofServerSetting 定义了 Pprof 服务器配置
var PprofServerSetting *httpServerConf

// WebServerSetting 定义了 Web 服务器配置
var WebServerSetting *httpServerConf

// AppSetting 定义了应用程序配置
var AppSetting *appConf

// CacheSetting 定义了缓存配置
var CacheSetting *cacheConf

// EventManagerSetting 定义了事件管理器配置
var EventManagerSetting *eventManagerConf

// MetricManagerSetting 定义了指标管理器配置
var MetricManagerSetting *metricManagerConf

// JobManagerSetting 定义了作业管理器配置
var JobManagerSetting *jobManagerConf

// CacheIndexSetting 定义了缓存索引配置
var CacheIndexSetting *cacheIndexConf

// SimpleCacheIndexSetting 定义了简单缓存索引配置
var SimpleCacheIndexSetting *simpleCacheIndexConf

// BigCacheIndexSetting 定义了大容量缓存索引配置
var BigCacheIndexSetting *bigCacheIndexConf

// RedisCacheIndexSetting 定义了 Redis 缓存索引配置
var RedisCacheIndexSetting *redisCacheIndexConf

// TweetSearchSetting 定义了推文搜索配置
var TweetSearchSetting *tweetSearchConf

// MeiliSetting 定义了 Meili 配置
var MeiliSetting *meiliConf

// ObjectStorage 定义了对象存储配置
var ObjectStorage *objectStorageConf

// MinIOSetting 定义了 MinIO 对象存储配置
var MinIOSetting *minioConf

// JWTSetting 定义了 JWT 配置
var JWTSetting *jwtConf

// WebProfileSetting 定义了 WebProfile 配置
var WebProfileSetting *WebProfileConf

// setupSetting 初始化配置参数
func setupSetting(suite []string, noDefault bool) error {
	vp, err := newViper()
	if err != nil {
		return err
	}

	// 初始化功能配置
	ss, kv := featuresInfoFrom(vp, "Features")
	cfg.Initial(ss, kv)
	if len(suite) > 0 {
		cfg.Use(suite, noDefault)
	}

	// 将各配置对象映射到变量
	objects := map[string]any{
		"App":              &AppSetting,
		"Cache":            &CacheSetting,
		"EventManager":     &EventManagerSetting,
		"MetricManager":    &MetricManagerSetting,
		"JobManager":       &JobManagerSetting,
		"PprofServer":      &PprofServerSetting,
		"WebServer":        &WebServerSetting,
		"CacheIndex":       &CacheIndexSetting,
		"SimpleCacheIndex": &SimpleCacheIndexSetting,
		"BigCacheIndex":    &BigCacheIndexSetting,
		"RedisCacheIndex":  &RedisCacheIndexSetting,
		"Pyroscope":        &PyroscopeSetting,
		"Logger":           &loggerSetting,
		"LoggerFile":       &loggerFileSetting,
		"LoggerMeili":      &loggerMeiliSetting,
		"Database":         &DatabaseSetting,
		"MySQL":            &MysqlSetting,
		"TweetSearch":      &TweetSearchSetting,
		"Meili":            &MeiliSetting,
		"Redis":            &redisSetting,
		"JWT":              &JWTSetting,
		"ObjectStorage":    &ObjectStorage,
		"MinIO":            &MinIOSetting,
		"WebProfile":       &WebProfileSetting,
	}

	// 解析配置并映射到各配置对象
	for k, v := range objects {
		err := vp.UnmarshalKey(k, v)
		if err != nil {
			return err
		}
	}

	// 设置一些配置参数的时间单位
	CacheSetting.CientSideCacheExpire *= time.Second
	EventManagerSetting.TickWaitTime *= time.Second
	MetricManagerSetting.TickWaitTime *= time.Second
	JWTSetting.Expire *= time.Second
	SimpleCacheIndexSetting.CheckTickDuration *= time.Second
	SimpleCacheIndexSetting.ExpireTickDuration *= time.Second
	BigCacheIndexSetting.ExpireInSecond *= time.Second
	RedisCacheIndexSetting.ExpireInSecond *= time.Second
	redisSetting.ConnWriteTimeout *= time.Second

	return nil
}

// Initial 初始化配置参数
func Initial(suite []string, noDefault bool) {
	err := setupSetting(suite, noDefault)
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	setupLogger()
}

// GetOssDomain 获取对象存储域名
func GetOssDomain() string {
	uri := "https://"
	if !MinIOSetting.Secure {
		uri = "http://"
	}
	return uri + MinIOSetting.Domain + "/" + MinIOSetting.Bucket + "/"
}

// RunMode 获取运行模式
func RunMode() string {
	return AppSetting.RunMode
}


// Copyright 2023 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package conf

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/pyroscope-io/client/pyroscope"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

//go:embed config.yaml
var configBytes []byte

// pyroscopeConf 定义 Pyroscope 配置结构体
type pyroscopeConf struct {
	AppName   string // 应用名称
	Endpoint  string // 终端点
	AuthToken string // 认证令牌
	Logger    string // 日志记录器
}



// loggerConf 定义 Logger 配置结构体
type loggerConf struct {
	Level string // 日志级别
}

type loggerFileConf struct {
	SavePath string
	FileName string
	FileExt  string
}

type loggerMeiliConf struct {
	Host         string
	Index        string
	ApiKey       string
	Secure       bool
	MaxLogBuffer int
	MinWorker    int
}

type httpServerConf struct {
	RunMode      string
	HttpIp       string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type appConf struct {
	RunMode               string
	MaxCommentCount       int64
	MaxWhisperDaily       int64
	MaxCaptchaTimes       int
	AttachmentIncomeRate  float64
	DefaultContextTimeout time.Duration
	DefaultPageSize       int
	MaxPageSize           int
}

type cacheConf struct {
	KeyPoolSize          int
	CientSideCacheExpire time.Duration
	UnreadMsgExpire      int64
	UserTweetsExpire     int64
	IndexTweetsExpire    int64
	MessagesExpire       int64
	IndexTrendsExpire    int64
	TweetCommentsExpire  int64
	OnlineUserExpire     int64
	UserInfoExpire       int64
	UserProfileExpire    int64
	UserRelationExpire   int64
}

type eventManagerConf struct {
	MinWorker       int
	MaxEventBuf     int
	MaxTempEventBuf int
	MaxTickCount    int
	TickWaitTime    time.Duration
}

type metricManagerConf struct {
	MinWorker       int
	MaxEventBuf     int
	MaxTempEventBuf int
	MaxTickCount    int
	TickWaitTime    time.Duration
}

type jobManagerConf struct {
	MaxOnlineInterval     string
	UpdateMetricsInterval string
}

type cacheIndexConf struct {
	MaxUpdateQPS int
	MinWorker    int
}

type simpleCacheIndexConf struct {
	MaxIndexSize       int
	CheckTickDuration  time.Duration
	ExpireTickDuration time.Duration
}

type bigCacheIndexConf struct {
	MaxIndexPage     int
	HardMaxCacheSize int
	ExpireInSecond   time.Duration
	Verbose          bool
}

type redisCacheIndexConf struct {
	ExpireInSecond time.Duration
	Verbose        bool
}

type tweetSearchConf struct {
	MaxUpdateQPS int
	MinWorker    int
}

type meiliConf struct {
	Host   string
	Index  string
	ApiKey string
	Secure bool
}

// databaseConf 定义数据库配置结构体
type databaseConf struct {
	TablePrefix string // 表前缀
	LogLevel    string // 日志级别
}

// mysqlConf 定义 MySQL 配置结构体
type mysqlConf struct {
	UserName     string // 用户名
	Password     string // 密码
	Host         string // 主机
	DBName       string // 数据库名
	Charset      string // 字符集
	ParseTime    bool   // 解析时间
	MaxIdleConns int    // 最大空闲连接数
	MaxOpenConns int    // 最大打开连接数
}

type objectStorageConf struct {
	RetainInDays int
	TempDir      string
}

type minioConf struct {
	AccessKey string
	SecretKey string
	Secure    bool
	Endpoint  string
	Bucket    string
	Domain    string
}

type redisConf struct {
	InitAddress      []string
	Username         string
	Password         string
	SelectDB         int
	ConnWriteTimeout time.Duration
}

type jwtConf struct {
	Secret string
	Issuer string
	Expire time.Duration
}

type WebProfileConf struct {
	UseFriendship           bool   `json:"use_friendship"`
	EnableTrendsBar         bool   `json:"enable_trends_bar"`
	EnableWallet            bool   `json:"enable_wallet"`
	AllowTweetAttachment    bool   `json:"allow_tweet_attachment"`
	AllowTweetVideo         bool   `json:"allow_tweet_video"`
	AllowUserRegister       bool   `json:"allow_user_register"`
	AllowPhoneBind          bool   `json:"allow_phone_bind"`
	DefaultTweetMaxLength   int    `json:"default_tweet_max_length"`
	TweetWebEllipsisSize    int    `json:"tweet_web_ellipsis_size"`
	TweetMobileEllipsisSize int    `json:"tweet_mobile_ellipsis_size"`
	DefaultTweetVisibility  string `json:"default_tweet_visibility"`
	DefaultMsgLoopInterval  int    `json:"default_msg_loop_interval"`
	CopyrightTop            string `json:"copyright_top"`
	CopyrightLeft           string `json:"copyright_left"`
	CopyrightLeftLink       string `json:"copyright_left_link"`
	CopyrightRight          string `json:"copyright_right"`
	CopyrightRightLink      string `json:"copyright_right_link"`
}

func (s *httpServerConf) GetReadTimeout() time.Duration {
	return s.ReadTimeout * time.Second
}

func (s *httpServerConf) GetWriteTimeout() time.Duration {
	return s.WriteTimeout * time.Second
}

func (s *mysqlConf) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		s.UserName,
		s.Password,
		s.Host,
		s.DBName,
		s.Charset,
		s.ParseTime,
	)
}

func (s *databaseConf) logLevel() logger.LogLevel {
	switch strings.ToLower(s.LogLevel) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Error
	}
}

func (s *databaseConf) TableNames() (res TableNameMap) {
	tableNames := []string{
		TableAnouncement,
		TableAnouncementContent,
		TableAttachment,
		TableCaptcha,
		TableComment,
		TableCommentMetric,
		TableCommentContent,
		TableCommentReply,
		TableFollowing,
		TableContact,
		TableContactGroup,
		TableMessage,
		TablePost,
		TablePostMetric,
		TablePostByComment,
		TablePostByMedia,
		TablePostAttachmentBill,
		TablePostCollection,
		TablePostContent,
		TablePostStar,
		TableTag,
		TableUser,
		TableUserRelation,
		TableUserMetric,
		TableWalletRecharge,
		TableWalletStatement,
	}
	res = make(TableNameMap, len(tableNames))
	for _, name := range tableNames {
		res[name] = s.TablePrefix + name
	}
	return
}

func (s *loggerConf) logLevel() logrus.Level {
	switch strings.ToLower(s.Level) {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.ErrorLevel
	}
}

func (s *loggerMeiliConf) Endpoint() string {
	return endpoint(s.Host, s.Secure)
}

func (s *loggerMeiliConf) minWork() int {
	if s.MinWorker < 5 {
		return 5
	} else if s.MinWorker > 100 {
		return 100
	}
	return s.MinWorker
}

func (s *loggerMeiliConf) maxLogBuffer() int {
	if s.MaxLogBuffer < 10 {
		return 10
	} else if s.MaxLogBuffer > 1000 {
		return 1000
	}
	return s.MaxLogBuffer
}

func (s *objectStorageConf) TempDirSlash() string {
	return strings.Trim(s.TempDir, " /") + "/"
}

func (s *meiliConf) Endpoint() string {
	return endpoint(s.Host, s.Secure)
}

func (s *pyroscopeConf) GetLogger() (logger pyroscope.Logger) {
	switch strings.ToLower(s.Logger) {
	case "standard":
		logger = pyroscope.StandardLogger
	case "logrus":
		logger = logrus.StandardLogger()
	}
	return
}

// endpoint 根据 host 和 secure 参数生成 URL 地址
func endpoint(host string, secure bool) string {
	schema := "http"
	if secure {
		schema = "https"
	}
	return schema + "://" + host
}

// newViper 返回一个配置了 config.yaml 的 Viper 实例
func newViper() (*viper.Viper, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath(".")
	vp.AddConfigPath("custom/")
	vp.SetConfigType("yaml")
	err := vp.ReadConfig(bytes.NewReader(configBytes))
	if err != nil {
		return nil, err
	}
	if err = vp.MergeInConfig(); err != nil {
		return nil, err
	}
	return vp, nil
}

// featuresInfoFrom 从 Viper 实例中提取特定键的信息
func featuresInfoFrom(vp *viper.Viper, k string) (map[string][]string, map[string]string) {
	sub := vp.Sub(k)
	keys := sub.AllKeys()

	suites := make(map[string][]string)
	kv := make(map[string]string, len(keys))
	for _, key := range sub.AllKeys() {
		val := sub.Get(key)
		switch v := val.(type) {
		case string:
			kv[key] = v
		case []any:
			suites[key] = sub.GetStringSlice(key)
		}
	}
	return suites, kv
}

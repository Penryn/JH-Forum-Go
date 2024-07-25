package conf

import (
	"sync"
	"time"

	"github.com/alimy/tryst/cfg"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var (
	_gormdb   *gorm.DB
	_onceGorm sync.Once
)

// MustGormDB 获取单例的 GORM 数据库连接对象
// 使用 sync.Once 确保在多线程环境下只初始化一次 GORM 数据库连接，并返回该连接对象。
func MustGormDB() *gorm.DB {
	_onceGorm.Do(func() {
		var err error
		if _gormdb, err = newGormDB(); err != nil {
			logrus.Fatalf("new gorm db failed: %s", err)
		}
	})
	return _gormdb
}

// newGormDB 创建一个新的 GORM 数据库连接
// 根据配置选择性地使用 MySQL 作为数据库，并配置日志、连接池等参数。
func newGormDB() (*gorm.DB, error) {
	newLogger := logger.New(
		logrus.StandardLogger(), // 日志输出的目标
		logger.Config{
			SlowThreshold:             time.Second,                // 慢查询阈值
			LogLevel:                  DatabaseSetting.logLevel(), // 日志级别
			IgnoreRecordNotFoundError: true,                       // 忽略记录未找到错误
			Colorful:                  false,                      // 禁用彩色打印
		},
	)

	config := &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   DatabaseSetting.TablePrefix, // 数据表前缀
			SingularTable: true,                        // 表名单数形式
		},
	}

	plugin := dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(time.Hour).              // 最大空闲连接时间
		SetConnMaxLifetime(24 * time.Hour).         // 最大连接持续时间
		SetMaxIdleConns(MysqlSetting.MaxIdleConns). // 最大空闲连接数
		SetMaxOpenConns(MysqlSetting.MaxOpenConns)  // 最大打开连接数

	var (
		db  *gorm.DB
		err error
	)
	if cfg.If("MySQL") {
		logrus.Debugln("use MySQL as db")
		if db, err = gorm.Open(mysql.Open(MysqlSetting.Dsn()), config); err == nil {
			db.Use(plugin)
		}
	} else {
		logrus.Debugln("use default of MySQL as db")
		if db, err = gorm.Open(mysql.Open(MysqlSetting.Dsn()), config); err == nil {
			db.Use(plugin)
		}
	}

	return db, err
}

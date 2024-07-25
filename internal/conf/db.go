package conf

import (
	"database/sql"
	"sync"

	"github.com/alimy/tryst/cfg"
	"github.com/sirupsen/logrus"
)

// 表名常量
const (
	TableAnouncement        = "user"
	TableAnouncementContent = "anouncement_content"
	TableAttachment         = "attachment"
	TableCaptcha            = "captcha"
	TableComment            = "comment"
	TableCommentMetric      = "comment_metric"
	TableCommentContent     = "comment_content"
	TableCommentReply       = "comment_reply"
	TableFollowing          = "following"
	TableContact            = "contact"
	TableContactGroup       = "contact_group"
	TableMessage            = "message"
	TablePost               = "post"
	TablePostMetric         = "post_metric"
	TablePostByComment      = "post_by_comment"
	TablePostByMedia        = "post_by_media"
	TablePostAttachmentBill = "post_attachment_bill"
	TablePostCollection     = "post_collection"
	TablePostContent        = "post_content"
	TablePostStar           = "post_star"
	TableTag                = "tag"
	TableUser               = "user"
	TableUserRelation       = "user_relation"
	TableUserMetric         = "user_metric"
	TableWalletRecharge     = "wallet_recharge"
	TableWalletStatement    = "wallet_statement"
)

type TableNameMap map[string]string

var (
	_sqldb   *sql.DB
	_onceSql sync.Once
)

// MustSqlDB 获取单例的 SQL 数据库连接对象
func MustSqlDB() *sql.DB {
	_onceSql.Do(func() {
		var err error
		if _, _sqldb, err = newSqlDB(); err != nil {
			logrus.Fatalf("new sql db failed: %s", err)
		}
	})
	return _sqldb
}

// newSqlDB 根据配置初始化 SQL 数据库连接对象
// 返回数据库驱动名称、数据库连接对象和可能的错误。
func newSqlDB() (driver string, db *sql.DB, err error) {
	if cfg.If("MySQL") {
		driver = "mysql"
		db, err = sql.Open(driver, MysqlSetting.Dsn())
	} else {
		driver = "mysql"
		db, err = sql.Open(driver, MysqlSetting.Dsn())
	}
	return
}

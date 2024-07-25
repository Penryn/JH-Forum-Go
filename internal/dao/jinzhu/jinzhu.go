// package jinzhu 实现了基于 gorm 的核心服务，支持多种数据库后端如 mysql。

package jinzhu

import (
	"sync"

	"github.com/Masterminds/semver/v3" // 引入 semver 库用于版本控制
	"JH-Forum/internal/conf"           // 引入配置包
	"JH-Forum/internal/core"           // 引入核心服务包
	"JH-Forum/internal/dao/cache"      // 引入缓存包
	"JH-Forum/internal/dao/security"   // 引入安全包
)

var (
	_onceInitial sync.Once // 用于确保初始化操作只执行一次的同步对象
)

type dataSrv struct {
	core.MessageService         // 消息服务接口
	core.TopicService           // 主题服务接口
	core.TweetService           // 推文服务接口
	core.TweetManageService     // 推文管理服务接口
	core.TweetHelpService       // 推文帮助服务接口
	core.TweetMetricServantA    // 推文指标服务接口A
	core.CommentService         // 评论服务接口
	core.CommentManageService   // 评论管理服务接口
	core.CommentMetricServantA  // 评论指标服务接口A
	core.TrendsManageServantA   // 趋势管理服务接口A
	core.UserManageService      // 用户管理服务接口
	core.UserMetricServantA     // 用户指标服务接口A
	core.ContactManageService   // 联系人管理服务接口
	core.FollowingManageService // 关注管理服务接口
	core.UserRelationService    // 用户关系服务接口
	core.SecurityService        // 安全服务接口
	core.AttachmentCheckService // 附件检查服务接口
}

type webDataSrvA struct {
	core.TopicServantA       // 主题服务接口A
	core.TweetServantA       // 推文服务接口A
	core.TweetManageServantA // 推文管理服务接口A
	core.TweetHelpServantA   // 推文帮助服务接口A
}

// NewDataService 返回一个包含数据服务接口和版本信息接口的实例
func NewDataService() (core.DataService, core.VersionInfo) {
	lazyInitial()                           // 执行懒初始化
	db := conf.MustGormDB()                 // 获取 GORM 数据库连接实例
	tms := newTweetMetricServentA(db)       // 实例化推文指标服务接口A
	ums := newUserMetricServentA(db)        // 实例化用户指标服务接口A
	cms := newCommentMetricServentA(db)     // 实例化评论指标服务接口A
	cis := cache.NewEventCacheIndexSrv(tms) // 实例化事件缓存索引服务
	ds := &dataSrv{                         // 创建数据服务实例
		TweetMetricServantA:    tms,
		CommentMetricServantA:  cms,
		UserMetricServantA:     ums,
		MessageService:         newMessageService(db),
		TopicService:           newTopicService(db),
		TweetService:           newTweetService(db),
		TweetManageService:     newTweetManageService(db, cis),
		TweetHelpService:       newTweetHelpService(db),
		CommentService:         newCommentService(db),
		CommentManageService:   newCommentManageService(db),
		TrendsManageServantA:   newTrendsManageServentA(db),
		UserManageService:      newUserManageService(db, ums),
		ContactManageService:   newContactManageService(db),
		FollowingManageService: newFollowingManageService(db),
		UserRelationService:    newUserRelationService(db),
		AttachmentCheckService: security.NewAttachmentCheckService(),
	}
	return cache.NewCacheDataService(ds), ds // 返回缓存数据服务实例和数据服务实例
}

// NewWebDataServantA 返回一个包含 Web 数据服务接口A和版本信息接口的实例
func NewWebDataServantA() (core.WebDataServantA, core.VersionInfo) {
	lazyInitial()           // 执行懒初始化
	db := conf.MustGormDB() // 获取 GORM 数据库连接实例
	ds := &webDataSrvA{     // 创建 Web 数据服务接口A 实例
		TopicServantA:       newTopicServantA(db),
		TweetServantA:       newTweetServantA(db),
		TweetManageServantA: newTweetManageServantA(db),
		TweetHelpServantA:   newTweetHelpServantA(db),
	}
	return ds, ds // 返回 Web 数据服务接口A 实例和版本信息接口实例
}

// NewAuthorizationManageService 返回授权管理服务实例
func NewAuthorizationManageService() core.AuthorizationManageService {
	return newAuthorizationManageService(conf.MustGormDB())
}

// Name 返回数据服务实例的名称
func (s *dataSrv) Name() string {
	return "Gorm"
}

// Version 返回数据服务实例的版本信息
func (s *dataSrv) Version() *semver.Version {
	return semver.MustParse("v0.2.0")
}

// Name 返回 Web 数据服务实例A 的名称
func (s *webDataSrvA) Name() string {
	return "Gorm"
}

// Version 返回 Web 数据服务实例A 的版本信息
func (s *webDataSrvA) Version() *semver.Version {
	return semver.MustParse("v0.1.0")
}

// lazyInitial 执行包的懒初始化操作，确保包内全局变量只初始化一次
func lazyInitial() {
	_onceInitial.Do(func() {
		initTableName() // 初始化数据表名
	})
}

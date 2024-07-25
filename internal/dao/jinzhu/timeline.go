// package jinzhu 实现了广场推文索引服务，包括根据用户ID查询广场推文列表等功能。

package jinzhu

import (
	"JH-Forum/internal/core"           // 引入核心服务包
	"JH-Forum/internal/core/cs"        // 引入核心服务-通用包
	"JH-Forum/internal/core/ms"        // 引入核心服务-消息服务包
	"JH-Forum/internal/dao/jinzhu/dbr" // 引入Jinzhu数据库访问包
	"JH-Forum/pkg/debug"               // 引入debug包
	"github.com/sirupsen/logrus"       // 引入logrus包
	"gorm.io/gorm"                     // 引入gorm包
)

var (
	_ core.IndexPostsService = (*shipIndexSrv)(nil)
	_ core.IndexPostsService = (*simpleIndexPostsSrv)(nil)
)

type shipIndexSrv struct {
	ams core.AuthorizationManageService
	ths core.TweetHelpService
	db  *gorm.DB
}

type simpleIndexPostsSrv struct {
	ths core.TweetHelpService
	db  *gorm.DB
}

// IndexPosts 根据userId查询广场推文列表，简单做到不同用户的主页都是不同的；
func (s *shipIndexSrv) IndexPosts(user *ms.User, offset int, limit int) (*ms.IndexTweetList, error) {
	predicates := dbr.Predicates{
		"ORDER": []any{"is_top DESC, latest_replied_on DESC"},
	}
	if user == nil {
		predicates["visibility = ?"] = []any{dbr.PostVisitPublic}
	} else if !user.IsAdmin {
		friendIds, _ := s.ams.BeFriendIds(user.ID)
		friendIds = append(friendIds, user.ID)
		args := []any{dbr.PostVisitPublic, dbr.PostVisitPrivate, user.ID, dbr.PostVisitFriend, friendIds}
		predicates["visibility = ? OR (visibility = ? AND user_id = ?) OR (visibility = ? AND user_id IN ?)"] = args
	}

	posts, err := (&dbr.Post{}).Fetch(s.db, predicates, offset, limit)
	if err != nil {
		logrus.Debugf("gormIndexPostsSrv.IndexPosts err: %v", err)
		return nil, err
	}
	formatPosts, err := s.ths.MergePosts(posts)
	if err != nil {
		return nil, err
	}

	total, err := (&dbr.Post{}).CountBy(s.db, predicates)
	if err != nil {
		return nil, err
	}

	return &ms.IndexTweetList{
		Tweets: formatPosts,
		Total:  total,
	}, nil
}

// simpleCacheIndexGetPosts simpleCacheIndex 专属获取广场推文列表函数
func (s *simpleIndexPostsSrv) IndexPosts(_user *ms.User, offset int, limit int) (*ms.IndexTweetList, error) {
	predicates := dbr.Predicates{
		"visibility = ?": []any{dbr.PostVisitPublic},
		"ORDER":          []any{"is_top DESC, latest_replied_on DESC"},
	}

	posts, err := (&dbr.Post{}).Fetch(s.db, predicates, offset, limit)
	if err != nil {
		logrus.Debugf("gormSimpleIndexPostsSrv.IndexPosts err: %v", err)
		return nil, err
	}

	formatPosts, err := s.ths.MergePosts(posts)
	if err != nil {
		return nil, err
	}

	total, err := (&dbr.Post{}).CountBy(s.db, predicates)
	if err != nil {
		return nil, err
	}

	return &ms.IndexTweetList{
		Tweets: formatPosts,
		Total:  total,
	}, nil
}

// TweetTimeline 获取用户的时间线推文列表，暂未实现
func (s *simpleIndexPostsSrv) TweetTimeline(userId int64, offset int, limit int) (*cs.TweetBox, error) {
	// TODO
	return nil, debug.ErrNotImplemented
}

// newShipIndexService 创建新的广场推文索引服务实例
func newShipIndexService(db *gorm.DB, ams core.AuthorizationManageService, ths core.TweetHelpService) core.IndexPostsService {
	return &shipIndexSrv{
		ams: ams,
		ths: ths,
		db:  db,
	}
}

// newSimpleIndexPostsService 创建新的简单广场推文索引服务实例
func newSimpleIndexPostsService(db *gorm.DB, ths core.TweetHelpService) core.IndexPostsService {
	return &simpleIndexPostsSrv{
		ths: ths,
		db:  db,
	}
}

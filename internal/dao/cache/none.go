package cache

import (
	"github.com/Masterminds/semver/v3"
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
	"JH-Forum/pkg/debug"
)

// noneCacheIndexServant 实现了不缓存索引服务的接口
type noneCacheIndexServant struct {
	ips core.IndexPostsService // 索引帖子服务接口
}

var (
	_ core.CacheIndexService = (*noneCacheIndexServant)(nil) // 确保实现了缓存索引服务接口
	_ core.VersionInfo       = (*noneCacheIndexServant)(nil) // 确保实现了版本信息接口
)

// IndexPosts 从索引服务获取帖子列表
func (s *noneCacheIndexServant) IndexPosts(user *ms.User, offset int, limit int) (*ms.IndexTweetList, error) {
	return s.ips.IndexPosts(user, offset, limit)
}

// TweetTimeline 获取用户的推文时间线，暂未实现
func (s *noneCacheIndexServant) TweetTimeline(userId int64, offset int, limit int) (*cs.TweetBox, error) {
	// TODO
	return nil, debug.ErrNotImplemented
}

// SendAction 发送操作到索引服务，空实现
func (s *noneCacheIndexServant) SendAction(_act core.IdxAct, _post *ms.Post) {
	// empty
}

// Name 返回服务名称
func (s *noneCacheIndexServant) Name() string {
	return "NoneCacheIndex"
}

// Version 返回服务版本
func (s *noneCacheIndexServant) Version() *semver.Version {
	return semver.MustParse("v0.1.0")
}

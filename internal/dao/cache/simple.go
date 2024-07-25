package cache

import (
	"sync/atomic"
	"time"

	"github.com/Masterminds/semver/v3"
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
	"JH-Forum/pkg/debug"
	"github.com/sirupsen/logrus"
)

var (
	_ core.CacheIndexService = (*simpleCacheIndexServant)(nil)
	_ core.VersionInfo       = (*simpleCacheIndexServant)(nil)
)

// simpleCacheIndexServant 实现了 CacheIndexService 和 VersionInfo 接口
type simpleCacheIndexServant struct {
	ips             core.IndexPostsService
	indexActionCh   chan core.IdxAct
	indexPosts      *ms.IndexTweetList
	atomicIndex     atomic.Value
	maxIndexSize    int
	checkTick       *time.Ticker
	expireIndexTick *time.Ticker
}

// IndexPosts 从缓存或数据库获取索引帖子
func (s *simpleCacheIndexServant) IndexPosts(user *ms.User, offset int, limit int) (*ms.IndexTweetList, error) {
	cacheResp := s.atomicIndex.Load().(*ms.IndexTweetList)
	end := offset + limit
	if cacheResp != nil {
		size := len(cacheResp.Tweets)
		logrus.Debugf("simpleCacheIndexServant.IndexPosts 从缓存获取帖子：%d，偏移：%d，限制：%d，开始：%d，结束：%d", size, offset, limit, offset, end)
		if size >= end {
			return &ms.IndexTweetList{
				Tweets: cacheResp.Tweets[offset:end],
				Total:  cacheResp.Total,
			}, nil
		}
	}

	logrus.Debugln("simpleCacheIndexServant.IndexPosts 从数据库获取帖子")
	return s.ips.IndexPosts(user, offset, limit)
}

// TweetTimeline 暂未实现
func (s *simpleCacheIndexServant) TweetTimeline(userId int64, offset int, limit int) (*cs.TweetBox, error) {
	// TODO
	return nil, debug.ErrNotImplemented
}

// SendAction 将索引相关操作发送到通道或协程中
func (s *simpleCacheIndexServant) SendAction(act core.IdxAct, _post *ms.Post) {
	select {
	case s.indexActionCh <- act:
		logrus.Debugf("simpleCacheIndexServant.SendAction 通过通道发送索引操作：%s", act)
	default:
		go func(ch chan<- core.IdxAct, act core.IdxAct) {
			logrus.Debugf("simpleCacheIndexServant.SendAction 通过协程发送索引操作：%s", act)
			ch <- act
		}(s.indexActionCh, act)
	}
}

// startIndexPosts 启动一个协程管理索引帖子
func (s *simpleCacheIndexServant) startIndexPosts() {
	var err error
	for {
		select {
		case <-s.checkTick.C:
			if s.indexPosts == nil {
				logrus.Debugf("通过 checkTick 索引帖子")
				if s.indexPosts, err = s.ips.IndexPosts(nil, 0, s.maxIndexSize); err == nil {
					s.atomicIndex.Store(s.indexPosts)
				} else {
					logrus.Errorf("获取索引帖子错误：%v", err)
				}
			}
		case <-s.expireIndexTick.C:
			logrus.Debugf("通过 expireIndexTick 过期索引帖子")
			if s.indexPosts != nil {
				s.indexPosts = nil
				s.atomicIndex.Store(s.indexPosts)
			}
		case action := <-s.indexActionCh:
			switch action {
			// TODO: 列出每种情况，以备后续精细化处理
			case core.IdxActCreatePost,
				core.IdxActUpdatePost,
				core.IdxActDeletePost,
				core.IdxActStickPost,
				core.IdxActVisiblePost:
				// 防止在最短时间内更新帖子太多
				if s.indexPosts != nil {
					logrus.Debugf("通过操作 %s 移除索引帖子", action)
					s.indexPosts = nil
					s.atomicIndex.Store(s.indexPosts)
				}
			default:
				// 空操作
			}
		}
	}
}

// Name 返回缓存实例的名称
func (s *simpleCacheIndexServant) Name() string {
	return "SimpleCacheIndex"
}

// Version 返回缓存实例的版本
func (s *simpleCacheIndexServant) Version() *semver.Version {
	return semver.MustParse("v0.2.0")
}

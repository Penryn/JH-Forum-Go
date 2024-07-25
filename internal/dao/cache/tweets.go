package cache

import (
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
	"github.com/sirupsen/logrus"
)

// eventCacheIndexSrv 实现了 CacheIndexService 接口，用于处理帖子索引事件
type eventCacheIndexSrv struct {
	tms core.TweetMetricServantA
}

// SendAction 根据索引操作发送相应事件，并更新帖子指标
func (s *eventCacheIndexSrv) SendAction(act core.IdxAct, post *ms.Post) {
	err := error(nil)
	switch act {
	case core.IdxActUpdatePost:
		err = s.tms.UpdateTweetMetric(&cs.TweetMetric{
			PostId:          post.ID,
			CommentCount:    post.CommentCount,
			UpvoteCount:     post.UpvoteCount,
			CollectionCount: post.CollectionCount,
			ShareCount:      post.ShareCount,
		})
		OnExpireIndexTweetEvent(post.UserID)
	case core.IdxActCreatePost:
		err = s.tms.AddTweetMetric(post.ID)
		OnExpireIndexTweetEvent(post.UserID)
	case core.IdxActDeletePost:
		err = s.tms.DeleteTweetMetric(post.ID)
		OnExpireIndexTweetEvent(post.UserID)
	case core.IdxActStickPost, core.IdxActVisiblePost:
		OnExpireIndexTweetEvent(post.UserID)
	}
	if err != nil {
		logrus.Errorf("eventCacheIndexSrv.SendAction(%s) 发生错误：%s", act, err)
	}
}

// NewEventCacheIndexSrv 创建一个新的 eventCacheIndexSrv 实例
func NewEventCacheIndexSrv(tms core.TweetMetricServantA) core.CacheIndexService {
	lazyInitial()
	return &eventCacheIndexSrv{
		tms: tms,
	}
}

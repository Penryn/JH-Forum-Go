// 该代码定义了两个接口，用于实现广场首页推文列表服务。

package core

import (
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
)

// IndexPostsService 广场首页推文列表服务接口
type IndexPostsService interface {
	IndexPosts(user *ms.User, offset int, limit int) (*ms.IndexTweetList, error) // 获取首页推文列表
}

// IndexPostsServantA 广场首页推文列表服务(版本A)接口
type IndexPostsServantA interface {
	IndexPosts(user *ms.User, limit int, offset int) (*cs.TweetBox, error) // 获取首页推文列表（版本A）
}

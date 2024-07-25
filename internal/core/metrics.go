// 该代码定义了三个接口，用于处理推文、评论和用户的度量更新操作。

package core

import (
	"JH-Forum/internal/core/cs"
)

// 推文度量服务接口
type TweetMetricServantA interface {
	UpdateTweetMetric(metric *cs.TweetMetric) error // 更新推文度量
	AddTweetMetric(postId int64) error              // 添加推文度量
	DeleteTweetMetric(postId int64) error           // 删除推文度量
}

// 评论度量服务接口
type CommentMetricServantA interface {
	UpdateCommentMetric(metric *cs.CommentMetric) error // 更新评论度量
	AddCommentMetric(commentId int64) error             // 添加评论度量
	DeleteCommentMetric(commentId int64) error          // 删除评论度量
}

// 用户度量服务接口
type UserMetricServantA interface {
	UpdateUserMetric(userId int64, action uint8) error // 更新用户度量
	AddUserMetric(userId int64) error                  // 添加用户度量
	DeleteUserMetric(userId int64) error               // 删除用户度量
}

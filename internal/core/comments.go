// 该文件定义了评论服务和评论管理服务的接口

package core

import (
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
)

// CommentService 评论检索服务接口
type CommentService interface {
	// 获取评论
	GetComments(tweetId int64, style cs.StyleCommentType, limit int, offset int) ([]*ms.Comment, int64, error)
	// 根据ID获取评论
	GetCommentByID(id int64) (*ms.Comment, error)
	// 根据ID获取评论回复
	GetCommentReplyByID(id int64) (*ms.CommentReply, error)
	// 根据ID获取评论内容
	GetCommentContentsByIDs(ids []int64) ([]*ms.CommentContent, error)
	// 根据ID获取评论回复内容
	GetCommentRepliesByID(ids []int64) ([]*ms.CommentReplyFormated, error)
	// 获取评论点赞状态
	GetCommentThumbsMap(userId int64, tweetId int64) (cs.CommentThumbsMap, cs.CommentThumbsMap, error)
}

// CommentManageService 评论管理服务接口
type CommentManageService interface {
	// 高亮评论
	HighlightComment(userId, commentId int64) (int8, error)
	// 删除评论
	DeleteComment(comment *ms.Comment) error
	// 创建评论
	CreateComment(comment *ms.Comment) (*ms.Comment, error)
	// 创建评论回复
	CreateCommentReply(reply *ms.CommentReply) (*ms.CommentReply, error)
	// 删除评论回复
	DeleteCommentReply(reply *ms.CommentReply) error
	// 创建评论内容
	CreateCommentContent(content *ms.CommentContent) (*ms.CommentContent, error)
	// 点赞评论
	ThumbsUpComment(userId int64, tweetId, commentId int64) error
	// 取消点赞评论
	ThumbsDownComment(userId int64, tweetId, commentId int64) error
	// 点赞评论回复
	ThumbsUpReply(userId int64, tweetId, commentId, replyId int64) error
	// 取消点赞评论回复
	ThumbsDownReply(userId int64, tweetId, commentId, replyId int64) error
}

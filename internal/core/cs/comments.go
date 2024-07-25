// 该代码定义了一些常量和结构体类型。

package cs

// StyleCommentType 表示评论样式类型
type StyleCommentType uint8

const (
	// StyleCommentDefault 默认评论样式
	StyleCommentDefault StyleCommentType = iota
	// StyleCommentHots 热门评论样式
	StyleCommentHots
	// StyleCommentNewest 最新评论样式
	StyleCommentNewest
)

// CommentThumbs 表示评论点赞结构体
type CommentThumbs struct {
	UserID       int64 `json:"user_id"`        // 用户ID
	TweetID      int64 `json:"tweet_id"`       // 推文ID
	CommentID    int64 `json:"comment_id"`     // 评论ID
	ReplyID      int64 `json:"reply_id"`       // 回复ID
	CommentType  int8  `json:"comment_type"`   // 评论类型
	IsThumbsUp   int8  `json:"is_thumbs_up"`   // 是否点赞
	IsThumbsDown int8  `json:"is_thumbs_down"` // 是否点踩
}

// CommentThumbsList 是 CommentThumbs 的切片类型
type CommentThumbsList []*CommentThumbs

// CommentThumbsMap 是 CommentThumbs 的映射类型
type CommentThumbsMap map[int64]*CommentThumbs

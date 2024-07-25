package dbr

// TweetCommentThumbs 结构体表示用户对动态评论或回复的点赞或踩操作记录。
type TweetCommentThumbs struct {
	*Model             // 模型基类
	UserID       int64 `json:"user_id"`        // 用户ID
	TweetID      int64 `json:"tweet_id"`       // 动态ID
	CommentID    int64 `json:"comment_id"`     // 评论ID
	ReplyID      int64 `json:"reply_id"`       // 回复ID
	CommentType  int8  `json:"comment_type"`   // 评论类型
	IsThumbsUp   int8  `json:"is_thumbs_up"`   // 是否点赞
	IsThumbsDown int8  `json:"is_thumbs_down"` // 是否踩
}

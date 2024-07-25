package cs

const (
	MetricActionCreateTweet uint8 = iota
	MetricActionDeleteTweet
)

// TweetMetric 表示推文的指标结构
type TweetMetric struct {
	PostId          int64 // 推文ID
	CommentCount    int64 // 评论数
	UpvoteCount     int64 // 点赞数
	CollectionCount int64 // 收藏数
	ShareCount      int64 // 分享数
	ThumbsUpCount   int64 // 点赞数
	ThumbsDownCount int64 // 点踩数
}

// CommentMetric 表示评论的指标结构
type CommentMetric struct {
	CommentId       int64 // 评论ID
	ReplyCount      int32 // 回复数
	ThumbsUpCount   int32 // 点赞数
	ThumbsDownCount int32 // 点踩数
}

// RankScore 计算推文的排名分数
func (m *TweetMetric) RankScore(motivationFactor int) int64 {
	if motivationFactor == 0 {
		motivationFactor = 1
	}
	return (m.CommentCount + m.UpvoteCount*2 + m.CollectionCount*4 + m.ShareCount*8) * int64(motivationFactor)
}

// RankScore 计算评论的排名分数
func (m *CommentMetric) RankScore(motivationFactor int) int64 {
	if motivationFactor == 0 {
		motivationFactor = 1
	}
	return int64(m.ReplyCount*2+m.ThumbsUpCount*4-m.ThumbsDownCount) * int64(motivationFactor)
}

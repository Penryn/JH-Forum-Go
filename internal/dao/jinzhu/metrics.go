// package jinzhu 提供了Tweet、Comment和User的指标服务实现，支持更新、添加和删除指标。

package jinzhu

import (
	"time"

	"JH-Forum/internal/core"           // 引入核心服务包
	"JH-Forum/internal/core/cs"        // 引入核心服务-客户端服务包
	"JH-Forum/internal/dao/jinzhu/dbr" // 引入Jinzhu数据库访问包
	"gorm.io/gorm"                     // 引入gorm包
)

type tweetMetricSrvA struct {
	db *gorm.DB // 数据库连接
}

type commentMetricSrvA struct {
	db *gorm.DB // 数据库连接
}

type userMetricSrvA struct {
	db *gorm.DB // 数据库连接
}

// UpdateTweetMetric 更新动态指标
func (s *tweetMetricSrvA) UpdateTweetMetric(metric *cs.TweetMetric) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		postMetric := &dbr.PostMetric{PostId: metric.PostId}
		tx.Model(postMetric).Where("post_id=?", metric.PostId).First(postMetric)
		postMetric.RankScore = metric.RankScore(postMetric.MotivationFactor)
		return tx.Save(postMetric).Error
	})
}

// AddTweetMetric 添加动态指标
func (s *tweetMetricSrvA) AddTweetMetric(postId int64) (err error) {
	_, err = (&dbr.PostMetric{PostId: postId}).Create(s.db)
	return
}

// DeleteTweetMetric 删除动态指标
func (s *tweetMetricSrvA) DeleteTweetMetric(postId int64) (err error) {
	return (&dbr.PostMetric{PostId: postId}).Delete(s.db)
}

// UpdateCommentMetric 更新评论指标
func (s *commentMetricSrvA) UpdateCommentMetric(metric *cs.CommentMetric) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		commentMetric := &dbr.CommentMetric{CommentId: metric.CommentId}
		tx.Model(commentMetric).Where("comment_id=?", metric.CommentId).First(commentMetric)
		commentMetric.RankScore = metric.RankScore(commentMetric.MotivationFactor)
		return tx.Save(commentMetric).Error
	})
}

// AddCommentMetric 添加评论指标
func (s *commentMetricSrvA) AddCommentMetric(commentId int64) (err error) {
	_, err = (&dbr.CommentMetric{CommentId: commentId}).Create(s.db)
	return
}

// DeleteCommentMetric 删除评论指标
func (s *commentMetricSrvA) DeleteCommentMetric(commentId int64) (err error) {
	return (&dbr.CommentMetric{CommentId: commentId}).Delete(s.db)
}

// UpdateUserMetric 更新用户指标
func (s *userMetricSrvA) UpdateUserMetric(userId int64, action uint8) error {
	metric := &dbr.UserMetric{UserId: userId}
	s.db.Model(metric).Where("user_id=?", userId).First(metric)
	metric.LatestTrendsOn = time.Now().Unix()
	switch action {
	case cs.MetricActionCreateTweet:
		metric.TweetsCount++
	case cs.MetricActionDeleteTweet:
		if metric.TweetsCount > 0 {
			metric.TweetsCount--
		}
	}
	return s.db.Save(metric).Error
}

// AddUserMetric 添加用户指标
func (s *userMetricSrvA) AddUserMetric(userId int64) (err error) {
	_, err = (&dbr.UserMetric{UserId: userId}).Create(s.db)
	return
}

// DeleteUserMetric 删除用户指标
func (s *userMetricSrvA) DeleteUserMetric(userId int64) (err error) {
	return (&dbr.UserMetric{UserId: userId}).Delete(s.db)
}

// newTweetMetricServentA 创建新的动态指标服务实例
func newTweetMetricServentA(db *gorm.DB) core.TweetMetricServantA {
	return &tweetMetricSrvA{
		db: db,
	}
}

// newCommentMetricServentA 创建新的评论指标服务实例
func newCommentMetricServentA(db *gorm.DB) core.CommentMetricServantA {
	return &commentMetricSrvA{
		db: db,
	}
}

// newUserMetricServentA 创建新的用户指标服务实例
func newUserMetricServentA(db *gorm.DB) core.UserMetricServantA {
	return &userMetricSrvA{
		db: db,
	}
}

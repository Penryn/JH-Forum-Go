// Copyright 2023 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dbr

import (
	"time"

	"gorm.io/gorm"
)

// PostMetric 表示帖子指标对象，包括帖子ID、排名分数、激励分数、衰减因子和动机因子等字段。
type PostMetric struct {
	*Model
	PostId           int64 `json:"post_id"`           // 帖子ID
	RankScore        int64 `json:"rank_score"`        // 排名分数
	IncentiveScore   int   `json:"incentive_score"`   // 激励分数
	DecayFactor      int   `json:"decay_factor"`      // 衰减因子
	MotivationFactor int   `json:"motivation_factor"` // 动机因子
}

// CommentMetric 表示评论指标对象，包括评论ID、排名分数、激励分数、衰减因子和动机因子等字段。
type CommentMetric struct {
	*Model
	CommentId        int64 `json:"comment_id"`        // 评论ID
	RankScore        int64 `json:"rank_score"`        // 排名分数
	IncentiveScore   int   `json:"incentive_score"`   // 激励分数
	DecayFactor      int   `json:"decay_factor"`      // 衰减因子
	MotivationFactor int   `json:"motivation_factor"` // 动机因子
}

// UserMetric 表示用户指标对象，包括用户ID、推文数量、最新趋势时间戳等字段。
type UserMetric struct {
	*Model
	UserId         int64 `json:"user_id"`          // 用户ID
	TweetsCount    int   `json:"tweets_count"`     // 推文数量
	LatestTrendsOn int64 `json:"latest_trends_on"` // 最新趋势时间戳
}

// Create 方法在数据库中创建新的帖子指标记录。
func (m *PostMetric) Create(db *gorm.DB) (*PostMetric, error) {
	err := db.Create(&m).Error
	return m, err
}

// Delete 方法在数据库中逻辑删除指定的帖子指标记录。
func (m *PostMetric) Delete(db *gorm.DB) error {
	return db.Model(m).Where("post_id", m.PostId).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// Create 方法在数据库中创建新的评论指标记录。
func (m *CommentMetric) Create(db *gorm.DB) (*CommentMetric, error) {
	err := db.Create(&m).Error
	return m, err
}

// Delete 方法在数据库中逻辑删除指定的评论指标记录。
func (m *CommentMetric) Delete(db *gorm.DB) error {
	return db.Model(m).Where("comment_id", m.CommentId).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// Create 方法在数据库中创建新的用户指标记录。
func (m *UserMetric) Create(db *gorm.DB) (*UserMetric, error) {
	err := db.Create(&m).Error
	return m, err
}

// Delete 方法在数据库中逻辑删除指定的用户指标记录。
func (m *UserMetric) Delete(db *gorm.DB) error {
	return db.Model(m).Where("user_id", m.UserId).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

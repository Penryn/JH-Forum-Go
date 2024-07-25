// Copyright 2022 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dbr

import (
	"time"

	"JH-Forum/pkg/types"
	"gorm.io/gorm"
)

// CommentReply 结构体表示评论回复对象，包括评论ID、用户ID、@用户ID、内容、IP地址、IP位置、点赞数和踩数。
type CommentReply struct {
	*Model                 // 模型基类
	CommentID       int64  `db:"comment_id" json:"comment_id"` // 评论ID
	UserID          int64  `db:"user_id" json:"user_id"`       // 用户ID
	AtUserID        int64  `db:"at_user_id" json:"at_user_id"` // @用户ID
	Content         string `json:"content"`                    // 内容
	IP              string `json:"ip"`                         // IP地址
	IPLoc           string `json:"ip_loc"`                     // IP位置
	ThumbsUpCount   int32  `json:"thumbs_up_count"`            // 点赞数
	ThumbsDownCount int32  `json:"-"`                          // 踩数
}

// CommentReplyFormated 结构体表示格式化后的评论回复对象，用于前端展示。
type CommentReplyFormated struct {
	ID            int64         `json:"id"`                         // ID
	CommentID     int64         `db:"comment_id" json:"comment_id"` // 评论ID
	UserID        int64         `db:"user_id" json:"user_id"`       // 用户ID
	User          *UserFormated `json:"user"`                       // 用户信息
	AtUserID      int64         `db:"at_user_id" json:"at_user_id"` // @用户ID
	AtUser        *UserFormated `json:"at_user"`                    // @用户信息
	Content       string        `json:"content"`                    // 内容
	IPLoc         string        `json:"ip_loc"`                     // IP位置
	ThumbsUpCount int32         `json:"thumbs_up_count"`            // 点赞数
	IsThumbsUp    int8          `json:"is_thumbs_up"`               // 是否点赞
	IsThumbsDown  int8          `json:"is_thumbs_down"`             // 是否踩
	CreatedOn     int64         `json:"created_on"`                 // 创建时间
	ModifiedOn    int64         `json:"modified_on"`                // 修改时间
}

// Format 将 CommentReply 对象格式化为 CommentReplyFormated 对象。
func (c *CommentReply) Format() *CommentReplyFormated {
	if c.Model == nil {
		return &CommentReplyFormated{}
	}

	return &CommentReplyFormated{
		ID:            c.ID,
		CommentID:     c.CommentID,
		UserID:        c.UserID,
		User:          &UserFormated{},
		AtUserID:      c.AtUserID,
		AtUser:        &UserFormated{},
		Content:       c.Content,
		IPLoc:         c.IPLoc,
		ThumbsUpCount: c.ThumbsUpCount,
		IsThumbsUp:    types.No,
		IsThumbsDown:  types.No,
		CreatedOn:     c.CreatedOn,
		ModifiedOn:    c.ModifiedOn,
	}
}

// List 根据条件从数据库中获取评论回复列表。
// 可以设置偏移量和限制数量来分页查询，还可以根据评论ID和其他条件进行过滤。
// 返回获取到的评论回复列表和可能发生的错误。
func (c *CommentReply) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*CommentReply, error) {
	var comments []*CommentReply
	var err error
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if c.CommentID > 0 {
		db = db.Where("id = ?", c.CommentID)
	}

	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}

	if err = db.Where("is_del = ?", 0).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

// Create 在数据库中创建评论回复记录。
// 返回创建的评论回复对象和可能发生的错误。
func (c *CommentReply) Create(db *gorm.DB) (*CommentReply, error) {
	err := db.Create(&c).Error
	return c, err
}

// Get 根据条件从数据库中获取单个评论回复记录。
// 返回获取到的评论回复对象和可能发生的错误。
func (c *CommentReply) Get(db *gorm.DB) (*CommentReply, error) {
	var reply CommentReply
	if c.Model != nil && c.ID > 0 {
		db = db.Where("id = ? AND is_del = ?", c.ID, 0)
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	err := db.First(&reply).Error
	if err != nil {
		return &reply, err
	}

	return &reply, nil
}

// Delete 在数据库中删除单个评论回复记录。
// 将删除标记设置为已删除，并记录删除时间。
// 返回可能发生的错误。
func (c *CommentReply) Delete(db *gorm.DB) error {
	return db.Model(&CommentReply{}).Where("id = ? AND is_del = ?", c.Model.ID, 0).Updates(map[string]any{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// DeleteByCommentIds 根据评论ID列表在数据库中删除相关的评论回复记录。
// 将删除标记设置为已删除，并记录删除时间。
// 返回可能发生的错误。
func (c *CommentReply) DeleteByCommentIds(db *gorm.DB, commentIds []int64) error {
	return db.Model(c).Where("comment_id IN ?", commentIds).Updates(map[string]any{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

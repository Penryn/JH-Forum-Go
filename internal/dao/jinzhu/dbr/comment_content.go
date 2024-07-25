package dbr

import (
	"time"

	"gorm.io/gorm"
)

// CommentContent 结构体表示评论内容对象，包括评论ID、用户ID、内容、类型和排序字段。
type CommentContent struct {
	*Model                 // 模型基类
	CommentID int64        `json:"comment_id"` // 评论ID
	UserID    int64        `json:"user_id"`    // 用户ID
	Content   string       `json:"content"`    // 内容
	Type      PostContentT `json:"type"`       // 类型
	Sort      int64        `json:"sort"`       // 排序字段
}

// List 根据条件从数据库中获取评论内容列表。
// 可以设置偏移量和限制数量来分页查询，还可以根据评论ID和其他条件进行过滤。
// 返回获取到的评论内容列表和可能发生的错误。
func (c *CommentContent) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*CommentContent, error) {
	var comments []*CommentContent
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

// Create 在数据库中创建评论内容记录。
// 返回创建的评论内容对象和可能发生的错误。
func (c *CommentContent) Create(db *gorm.DB) (*CommentContent, error) {
	err := db.Create(&c).Error
	return c, err
}

// MediaContentsByCommentId 根据评论ID列表从数据库中获取媒体内容。
// 只获取类型为图片的内容。
// 返回获取到的媒体内容列表和可能发生的错误。
func (c *CommentContent) MediaContentsByCommentId(db *gorm.DB, commentIds []int64) (contents []string, err error) {
	err = db.Model(c).Where("comment_id IN ? AND type = ?", commentIds, ContentTypeImage).Select("content").Find(&contents).Error
	return
}

// DeleteByCommentIds 根据评论ID列表在数据库中删除评论内容。
// 将删除标记设置为已删除，并记录删除时间。
// 返回可能发生的错误。
func (c *CommentContent) DeleteByCommentIds(db *gorm.DB, commentIds []int64) error {
	return db.Model(c).Where("comment_id IN ?", commentIds).Updates(map[string]any{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

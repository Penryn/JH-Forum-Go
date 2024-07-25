package dbr

import (
	"time"

	"JH-Forum/pkg/types"
	"gorm.io/gorm"
)

// Comment 表示系统中的评论对象，包括帖子ID、用户ID、IP地址、IP位置、精华标志、回复数、点赞数和踩数。
type Comment struct {
	*Model                 // 模型基类
	PostID          int64  `json:"post_id"`         // 帖子ID
	UserID          int64  `json:"user_id"`         // 用户ID
	IP              string `json:"ip"`              // IP地址
	IPLoc           string `json:"ip_loc"`          // IP位置
	IsEssence       int8   `json:"is_essense"`      // 是否精华
	ReplyCount      int32  `json:"reply_count"`     // 回复数
	ThumbsUpCount   int32  `json:"thumbs_up_count"` // 点赞数
	ThumbsDownCount int32  `json:"-"`               // 踩数
}

// CommentFormated 表示格式化后的评论对象，包括ID、帖子ID、用户ID、用户信息、评论内容、回复列表、IP位置、回复数、点赞数、是否精华等信息。
type CommentFormated struct {
	ID            int64                   `json:"id"`              // ID
	PostID        int64                   `json:"post_id"`         // 帖子ID
	UserID        int64                   `json:"user_id"`         // 用户ID
	User          *UserFormated           `json:"user"`            // 用户信息
	Contents      []*CommentContent       `json:"contents"`        // 评论内容列表
	Replies       []*CommentReplyFormated `json:"replies"`         // 回复列表
	IPLoc         string                  `json:"ip_loc"`          // IP位置
	ReplyCount    int32                   `json:"reply_count"`     // 回复数
	ThumbsUpCount int32                   `json:"thumbs_up_count"` // 点赞数
	IsEssence     int8                    `json:"is_essence"`      // 是否精华
	IsThumbsUp    int8                    `json:"is_thumbs_up"`    // 是否点赞
	IsThumbsDown  int8                    `json:"is_thumbs_down"`  // 是否踩
	CreatedOn     int64                   `json:"created_on"`      // 创建时间
	ModifiedOn    int64                   `json:"modified_on"`     // 修改时间
}

// Format 方法将 Comment 结构体格式化为 CommentFormated 结构体。
func (c *Comment) Format() *CommentFormated {
	if c.Model == nil {
		return &CommentFormated{}
	}
	return &CommentFormated{
		ID:            c.Model.ID,
		PostID:        c.PostID,
		UserID:        c.UserID,
		User:          &UserFormated{},
		Contents:      []*CommentContent{},
		Replies:       []*CommentReplyFormated{},
		IPLoc:         c.IPLoc,
		ReplyCount:    c.ReplyCount,
		ThumbsUpCount: c.ThumbsUpCount,
		IsEssence:     c.IsEssence,
		IsThumbsUp:    types.No,
		IsThumbsDown:  types.No,
		CreatedOn:     c.CreatedOn,
		ModifiedOn:    c.ModifiedOn,
	}
}

// Get 方法从数据库中获取指定的评论记录。
func (c *Comment) Get(db *gorm.DB) (*Comment, error) {
	var comment Comment
	if c.Model != nil && c.ID > 0 {
		db = db.Where("id = ? AND is_del = ?", c.ID, 0)
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	err := db.First(&comment).Error
	if err != nil {
		return &comment, err
	}

	return &comment, nil
}

// List 方法从数据库中按条件获取评论列表。
func (c *Comment) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*Comment, error) {
	var comments []*Comment
	var err error
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if c.PostID > 0 {
		db = db.Where("post_id = ?", c.PostID)
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

// Count 方法统计符合条件的评论数量。
func (c *Comment) Count(db *gorm.DB, conditions *ConditionsT) (int64, error) {
	var count int64
	if c.PostID > 0 {
		db = db.Where("post_id = ?", c.PostID)
	}
	for k, v := range *conditions {
		if k != "ORDER" {
			db = db.Where(k, v)
		}
	}
	if err := db.Model(c).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Create 方法在数据库中创建新的评论记录。
func (c *Comment) Create(db *gorm.DB) (*Comment, error) {
	err := db.Create(&c).Error
	return c, err
}

// Delete 方法从数据库中删除指定评论记录。
func (c *Comment) Delete(db *gorm.DB) error {
	return db.Model(c).Where("id = ?", c.Model.ID).Updates(map[string]any{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// CommentIdsByPostId 方法根据帖子ID获取相关评论的ID列表。
func (c *Comment) CommentIdsByPostId(db *gorm.DB, postId int64) (ids []int64, err error) {
	err = db.Model(c).Where("post_id = ?", postId).Select("id").Find(&ids).Error
	return
}

// DeleteByPostId 方法根据帖子ID删除相关评论。
func (c *Comment) DeleteByPostId(db *gorm.DB, postId int64) error {
	return db.Model(c).Where("post_id = ?", postId).Updates(map[string]any{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// PostContent 定义了帖子内容模型，包含帖子ID、用户ID、内容、类型和排序等信息。
package dbr

import (
	"time"

	"gorm.io/gorm"
)

// PostContentT 表示帖子内容的类型。
type PostContentT int

// 定义了帖子内容的类型常量。
const (
	ContentTypeTitle PostContentT = iota + 1
	ContentTypeText
	ContentTypeImage
	ContentTypeVideo
	ContentTypeAudio
	ContentTypeLink
	ContentTypeAttachment
	ContentTypeChargeAttachment
)

// mediaContentType 存储了所有媒体类型的常量。
var (
	mediaContentType = []PostContentT{
		ContentTypeImage,
		ContentTypeVideo,
		ContentTypeAudio,
		ContentTypeAttachment,
		ContentTypeChargeAttachment,
	}
)

// PostContent 表示帖子内容的结构。
type PostContent struct {
	*Model
	PostID  int64        `json:"post_id"` // 帖子ID
	UserID  int64        `json:"user_id"` // 用户ID
	Content string       `json:"content"` // 内容
	Type    PostContentT `json:"type"`    // 类型
	Sort    int64        `json:"sort"`    // 排序
}

// PostContentFormated 表示格式化后的帖子内容结构。
type PostContentFormated struct {
	ID      int64        `db:"id" json:"id"` // ID
	PostID  int64        `json:"post_id"`    // 帖子ID
	Content string       `json:"content"`    // 内容
	Type    PostContentT `json:"type"`       // 类型
	Sort    int64        `json:"sort"`       // 排序
}

// DeleteByPostId 根据帖子ID删除帖子内容。
func (p *PostContent) DeleteByPostId(db *gorm.DB, postId int64) error {
	return db.Model(p).Where("post_id = ?", postId).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// MediaContentsByPostId 根据帖子ID获取媒体内容。
func (p *PostContent) MediaContentsByPostId(db *gorm.DB, postId int64) (contents []string, err error) {
	err = db.Model(p).Where("post_id = ? AND type IN ?", postId, mediaContentType).Select("content").Find(&contents).Error
	return
}

// Create 创建帖子内容记录。
func (p *PostContent) Create(db *gorm.DB) (*PostContent, error) {
	err := db.Create(&p).Error
	return p, err
}

// Format 格式化帖子内容。
func (p *PostContent) Format() *PostContentFormated {
	if p.Model == nil {
		return nil
	}
	return &PostContentFormated{
		ID:      p.ID,
		PostID:  p.PostID,
		Content: p.Content,
		Type:    p.Type,
		Sort:    p.Sort,
	}
}

// List 根据条件查询帖子内容列表。
func (p *PostContent) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*PostContent, error) {
	var contents []*PostContent
	var err error
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if p.PostID > 0 {
		db = db.Where("id = ?", p.PostID)
	}

	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}

	if err = db.Where("is_del = ?", 0).Find(&contents).Error; err != nil {
		return nil, err
	}

	return contents, nil
}

// Get 根据条件获取单个帖子内容记录。
func (p *PostContent) Get(db *gorm.DB) (*PostContent, error) {
	var content PostContent
	if p.Model != nil && p.ID > 0 {
		db = db.Where("id = ? AND is_del = ?", p.ID, 0)
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	err := db.First(&content).Error
	if err != nil {
		return &content, err
	}

	return &content, nil
}

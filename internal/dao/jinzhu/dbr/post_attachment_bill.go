// PostAttachmentBill 定义了帖子附件账单模型，包含帖子ID、用户ID和支付金额等信息。
package dbr

import "gorm.io/gorm"

// PostAttachmentBill 表示帖子附件的账单信息。
type PostAttachmentBill struct {
	*Model
	PostID     int64 `json:"post_id"`     // 帖子ID
	UserID     int64 `json:"user_id"`     // 用户ID
	PaidAmount int64 `json:"paid_amount"` // 支付金额
}

// Get 根据条件获取单个帖子附件账单记录。
func (p *PostAttachmentBill) Get(db *gorm.DB) (*PostAttachmentBill, error) {
	var pas PostAttachmentBill
	if p.Model != nil && p.ID > 0 {
		db = db.Where("id = ? AND is_del = ?", p.ID, 0)
	}
	if p.PostID > 0 {
		db = db.Where("post_id = ?", p.PostID)
	}
	if p.UserID > 0 {
		db = db.Where("user_id = ?", p.UserID)
	}

	err := db.First(&pas).Error
	if err != nil {
		return &pas, err
	}

	return &pas, nil
}

// Create 创建帖子附件账单记录。
func (p *PostAttachmentBill) Create(db *gorm.DB) (*PostAttachmentBill, error) {
	err := db.Create(&p).Error
	return p, err
}

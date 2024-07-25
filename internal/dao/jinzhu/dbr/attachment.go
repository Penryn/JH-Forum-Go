package dbr

import "gorm.io/gorm"

// AttachmentType 表示附件的类型。
type AttachmentType int

const (
	AttachmentTypeImage AttachmentType = iota + 1
	AttachmentTypeVideo
	AttachmentTypeOther
)

// Attachment 是一个结构体，用于表示系统中的附件对象，包括用户ID、文件大小、图片宽度、高度、附件类型和内容。
type Attachment struct {
	*Model                   // 附件模型
	UserID    int64          `json:"user_id"`    // 用户ID
	FileSize  int64          `json:"file_size"`  // 文件大小
	ImgWidth  int            `json:"img_width"`  // 图片宽度
	ImgHeight int            `json:"img_height"` // 图片高度
	Type      AttachmentType `json:"type"`       // 附件类型
	Content   string         `json:"content"`    // 内容
}

// Create 在数据库中创建附件记录。
// 返回创建的附件对象和可能发生的错误。
func (a *Attachment) Create(db *gorm.DB) (*Attachment, error) {
	err := db.Create(&a).Error
	return a, err
}

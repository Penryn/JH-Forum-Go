package dbr

import "gorm.io/gorm"

type MessageT int8

const (
	MsgTypePost MessageT = iota + 1
	MsgtypeComment
	MsgTypeReply
	MsgTypeWhisper
	MsgTypeRequestingFriend
	MsgTypeSystem MessageT = 99

	MsgStatusUnread = 0
	MsgStatusReaded = 1
)

// Message 表示消息对象，包括发送者、接收者、消息类型、内容等字段。
type Message struct {
	*Model
	SenderUserID   int64    `json:"sender_user_id"`   // 发送者用户ID
	ReceiverUserID int64    `json:"receiver_user_id"` // 接收者用户ID
	Type           MessageT `json:"type"`             // 消息类型
	Brief          string   `json:"brief"`            // 简要
	Content        string   `json:"content"`          // 内容
	PostID         int64    `json:"post_id"`          // 帖子ID
	CommentID      int64    `json:"comment_id"`       // 评论ID
	ReplyID        int64    `json:"reply_id"`         // 回复ID
	IsRead         int8     `json:"is_read"`          // 是否已读
}

// MessageFormated 表示格式化后的消息对象，包含关联的用户、帖子、评论和回复等信息。
type MessageFormated struct {
	ID             int64         `json:"id"`                      // 消息ID
	SenderUserID   int64         `json:"sender_user_id"`          // 发送者用户ID
	SenderUser     *UserFormated `json:"sender_user"`             // 发送者用户信息
	ReceiverUserID int64         `json:"receiver_user_id"`        // 接收者用户ID
	ReceiverUser   *UserFormated `json:"receiver_user,omitempty"` // 接收者用户信息
	Type           MessageT      `json:"type"`                    // 消息类型
	Brief          string        `json:"brief"`                   // 简要
	Content        string        `json:"content"`                 // 内容
	PostID         int64         `json:"post_id"`                 // 帖子ID
	Post           *PostFormated `json:"post"`                    // 关联的帖子信息
	CommentID      int64         `json:"comment_id"`              // 评论ID
	Comment        *Comment      `json:"comment"`                 // 关联的评论信息
	ReplyID        int64         `json:"reply_id"`                // 回复ID
	Reply          *CommentReply `json:"reply"`                   // 关联的回复信息
	IsRead         int8          `json:"is_read"`                 // 是否已读
	CreatedOn      int64         `json:"created_on"`              // 创建时间戳
	ModifiedOn     int64         `json:"modified_on"`             // 修改时间戳
}

// Format 方法将 Message 对象格式化为 MessageFormated 对象。
func (m *Message) Format() *MessageFormated {
	if m.Model == nil || m.Model.ID == 0 {
		return nil
	}
	mf := &MessageFormated{
		ID:             m.ID,
		SenderUserID:   m.SenderUserID,
		SenderUser:     &UserFormated{},
		ReceiverUserID: m.ReceiverUserID,
		ReceiverUser:   &UserFormated{},
		Type:           m.Type,
		Brief:          m.Brief,
		Content:        m.Content,
		PostID:         m.PostID,
		Post:           &PostFormated{},
		CommentID:      m.CommentID,
		Comment:        &Comment{},
		ReplyID:        m.ReplyID,
		Reply:          &CommentReply{},
		IsRead:         m.IsRead,
		CreatedOn:      m.CreatedOn,
		ModifiedOn:     m.ModifiedOn,
	}

	return mf
}

// Create 方法在数据库中创建新的消息记录。
func (m *Message) Create(db *gorm.DB) (*Message, error) {
	err := db.Create(&m).Error
	return m, err
}

// Update 方法更新数据库中的消息记录。
func (m *Message) Update(db *gorm.DB) error {
	return db.Model(&Message{}).Where("id = ? AND is_del = ?", m.Model.ID, 0).Save(m).Error
}

// Get 方法从数据库中获取指定消息记录。
func (m *Message) Get(db *gorm.DB) (*Message, error) {
	var message Message
	if m.Model != nil && m.ID > 0 {
		db = db.Where("id = ? AND is_del = ?", m.ID, 0)
	}
	if m.ReceiverUserID > 0 {
		db = db.Where("receiver_user_id = ?", m.ReceiverUserID)
	}
	if err := db.First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

// FetchBy 方法从数据库中获取符合条件的消息记录。
func (m *Message) FetchBy(db *gorm.DB, predicates Predicates) ([]*Message, error) {
	var messages []*Message
	for k, v := range predicates {
		db = db.Where(k, v...)
	}
	db = db.Where("is_del = 0")
	if err := db.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// List 方法从数据库中获取指定用户接收的消息列表。
func (c *Message) List(db *gorm.DB, userId int64, offset, limit int) (res []*Message, err error) {
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	err = db.Where("receiver_user_id=? OR (sender_user_id=? AND type=4)", userId, userId).Order("id DESC").Find(&res).Error
	return
}

// Count 方法统计指定用户接收的消息数量。
func (m *Message) Count(db *gorm.DB, userId int64) (res int64, err error) {
	err = db.Model(m).Where("receiver_user_id=? OR (sender_user_id=? AND type=4)", userId, userId).Count(&res).Error
	return
}

// CountUnread 方法统计指定用户未读消息数量。
func (m *Message) CountUnread(db *gorm.DB, userId int64) (res int64, err error) {
	err = db.Model(m).Where("receiver_user_id=? AND is_read=0", userId).Count(&res).Error
	return
}

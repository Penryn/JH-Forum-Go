// package jinzhu 实现了消息服务接口，支持创建消息、获取未读消息数、按条件获取消息列表等功能。

package jinzhu

import (
	"JH-Forum/internal/core"           // 引入核心服务包
	"JH-Forum/internal/core/cs"        // 引入核心服务-客户端服务包
	"JH-Forum/internal/core/ms"        // 引入核心服务-消息服务包
	"JH-Forum/internal/dao/jinzhu/dbr" // 引入Jinzhu数据库访问包
	"gorm.io/gorm"                     // 引入gorm包
)

var (
	_ core.MessageService = (*messageSrv)(nil) // 确保messageSrv实现了MessageService接口
)

type messageSrv struct {
	db *gorm.DB // 数据库连接
}

// newMessageService 创建一个新的消息服务实例
func newMessageService(db *gorm.DB) core.MessageService {
	return &messageSrv{
		db: db,
	}
}

// CreateMessage 创建消息
func (s *messageSrv) CreateMessage(msg *ms.Message) (*ms.Message, error) {
	return msg.Create(s.db)
}

// GetUnreadCount 获取未读消息数
func (s *messageSrv) GetUnreadCount(userID int64) (int64, error) {
	return (&dbr.Message{}).CountUnread(s.db, userID)
}

// GetMessageByID 根据消息ID获取消息
func (s *messageSrv) GetMessageByID(id int64) (*ms.Message, error) {
	return (&dbr.Message{
		Model: &dbr.Model{
			ID: id,
		},
	}).Get(s.db)
}

// ReadMessage 标记消息为已读
func (s *messageSrv) ReadMessage(message *ms.Message) error {
	message.IsRead = 1
	return message.Update(s.db)
}

// ReadAllMessage 将用户的所有消息标记为已读
func (s *messageSrv) ReadAllMessage(userId int64) error {
	return s.db.Table(_message_).Where("receiver_user_id=? AND is_del=0", userId).Update("is_read", 1).Error
}

// GetMessages 根据条件获取消息列表
func (s *messageSrv) GetMessages(userId int64, style cs.MessageStyle, limit int, offset int) (res []*ms.MessageFormated, total int64, err error) {
	var messages []*dbr.Message
	db := s.db.Table(_message_)

	switch style {
	case cs.StyleMsgSystem:
		db = db.Where("receiver_user_id=? AND type IN (1, 2, 3, 99)", userId)
	case cs.StyleMsgWhisper:
		db = db.Where("(receiver_user_id=? OR sender_user_id=?) AND type=4", userId, userId)
	case cs.StyleMsgRequesting:
		db = db.Where("receiver_user_id=? AND type=5", userId)
	case cs.StyleMsgUnread:
		db = db.Where("receiver_user_id=? AND is_read=0", userId)
	case cs.StyleMsgAll:
		fallthrough
	default:
		db = db.Where("receiver_user_id=? OR (sender_user_id=? AND type=4)", userId, userId)
	}

	if err = db.Count(&total).Error; err != nil || total == 0 {
		return
	}

	if offset >= 0 && limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Order("id DESC").Find(&messages).Error; err != nil {
		return
	}

	for _, message := range messages {
		res = append(res, message.Format())
	}

	return
}

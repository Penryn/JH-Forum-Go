// 该代码定义了消息服务的接口。

package core

import (
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
)

// MessageService 消息服务接口
type MessageService interface {
	CreateMessage(msg *ms.Message) (*ms.Message, error)                                                       // 创建消息
	GetUnreadCount(userID int64) (int64, error)                                                               // 获取未读消息数量
	GetMessageByID(id int64) (*ms.Message, error)                                                             // 根据ID获取消息
	ReadMessage(message *ms.Message) error                                                                    // 标记消息为已读
	ReadAllMessage(userId int64) error                                                                        // 标记所有消息为已读
	GetMessages(userId int64, style cs.MessageStyle, limit, offset int) ([]*ms.MessageFormated, int64, error) // 获取消息列表
}

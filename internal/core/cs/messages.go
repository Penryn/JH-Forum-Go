// 该代码定义了消息列表的样式常量及类型。

package cs

// MessageStyle 表示消息列表的样式类型
type MessageStyle string

const (
	// 消息列表样式
	StyleMsgAll        MessageStyle = "all"        // 所有消息
	StyleMsgSystem     MessageStyle = "system"     // 系统消息
	StyleMsgWhisper    MessageStyle = "whisper"    // 悄悄话消息
	StyleMsgRequesting MessageStyle = "requesting" // 好友请求消息
	StyleMsgUnread     MessageStyle = "unread"     // 未读消息
)

package web

// AuditStyle 定义了审计样式的枚举类型。
type AuditStyle uint8

const (
	// AuditStyleUnknown 表示未知的审计样式。
	AuditStyleUnknown AuditStyle = iota
	// AuditStyleUserTweet 表示用户发帖审计样式。
	AuditStyleUserTweet
	// AuditStyleUserTweetComment 表示用户发帖评论审计样式。
	AuditStyleUserTweetComment
	// AuditStyleUserTweetReply 表示用户发帖回复审计样式。
	AuditStyleUserTweetReply
)

const (
	// AuditHookCtxKey 定义了审计上下文键。
	AuditHookCtxKey = "audit_ctx_key"
	// OnlineUserCtxKey 定义了在线用户上下文键。
	OnlineUserCtxKey = "online_user_ctx_key"
)

// AuditMetaInfo 定义了审计元信息结构体，包括审计样式和ID。
type AuditMetaInfo struct {
	Style AuditStyle
	Id    int64
}

// String 实现了 AuditStyle 类型的 String 方法，用于返回审计样式的字符串表示。
func (s AuditStyle) String() string {
	switch s {
	case AuditStyleUserTweet:
		return "UserTweet"
	case AuditStyleUserTweetComment:
		return "UserTweetComment"
	case AuditStyleUserTweetReply:
		return "UserTweetReply"
	case AuditStyleUnknown:
		fallthrough
	default:
		return "Unknown"
	}
}

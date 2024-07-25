// 该代码定义了一些常量和类型的别名，与外部包中的类型对应。

package ms

import (
	"JH-Forum/internal/dao/jinzhu/dbr"
)

const (
	MsgTypePost             = dbr.MsgTypePost             // MsgTypePost 的别名
	MsgtypeComment          = dbr.MsgtypeComment          // MsgtypeComment 的别名
	MsgTypeReply            = dbr.MsgTypeReply            // MsgTypeReply 的别名
	MsgTypeWhisper          = dbr.MsgTypeWhisper          // MsgTypeWhisper 的别名
	MsgTypeRequestingFriend = dbr.MsgTypeRequestingFriend // MsgTypeRequestingFriend 的别名
	MsgTypeSystem           = dbr.MsgTypeSystem           // MsgTypeSystem 的别名

	MsgStatusUnread = dbr.MsgStatusUnread // MsgStatusUnread 的别名
	MsgStatusReaded = dbr.MsgStatusReaded // MsgStatusReaded 的别名
)

type (
	Message         = dbr.Message         // Message 的别名
	MessageFormated = dbr.MessageFormated // MessageFormated 的别名
)

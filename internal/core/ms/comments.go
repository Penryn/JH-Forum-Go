// 该代码定义了一些类型的别名，与外部包中的类型对应。

package ms

import (
	"JH-Forum/internal/dao/jinzhu/dbr"
)

type (
	Comment              = dbr.Comment              // Comment 的别名
	CommentFormated      = dbr.CommentFormated      // CommentFormated 的别名
	CommentReply         = dbr.CommentReply         // CommentReply 的别名
	CommentContent       = dbr.CommentContent       // CommentContent 的别名
	CommentReplyFormated = dbr.CommentReplyFormated // CommentReplyFormated 的别名
)

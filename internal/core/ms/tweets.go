// 该代码定义了一些常量和类型的别名，与外部包中的类型对应。

package ms

import (
	"JH-Forum/internal/dao/jinzhu/dbr"
)

const (
	AttachmentTypeImage = dbr.AttachmentTypeImage // AttachmentTypeImage 的别名
	AttachmentTypeVideo = dbr.AttachmentTypeVideo // AttachmentTypeVideo 的别名
	AttachmentTypeOther = dbr.AttachmentTypeOther // AttachmentTypeOther 的别名

	// ContentTypeTitle 到 ContentTypeChargeAttachment 的常量别名
	ContentTypeTitle            = dbr.ContentTypeTitle
	ContentTypeText             = dbr.ContentTypeText
	ContentTypeImage            = dbr.ContentTypeImage
	ContentTypeVideo            = dbr.ContentTypeVideo
	ContentTypeAudio            = dbr.ContentTypeAudio
	ContentTypeLink             = dbr.ContentTypeLink
	ContentTypeAttachment       = dbr.ContentTypeAttachment
	ContentTypeChargeAttachment = dbr.ContentTypeChargeAttachment
)

const (
	PostVisitPublic    = dbr.PostVisitPublic    // PostVisitPublic 的别名
	PostVisitPrivate   = dbr.PostVisitPrivate   // PostVisitPrivate 的别名
	PostVisitFriend    = dbr.PostVisitFriend    // PostVisitFriend 的别名
	PostVisitFollowing = dbr.PostVisitFollowing // PostVisitFollowing 的别名
)

type (
	PostStar           = dbr.PostStar           // PostStar 的别名
	PostCollection     = dbr.PostCollection     // PostCollection 的别名
	PostAttachmentBill = dbr.PostAttachmentBill // PostAttachmentBill 的别名
	PostContent        = dbr.PostContent        // PostContent 的别名
	Attachment         = dbr.Attachment         // Attachment 的别名
	AttachmentType     = dbr.AttachmentType     // AttachmentType 的别名
	PostContentT       = dbr.PostContentT       // PostContentT 的别名
	PostVisibleT       = dbr.PostVisibleT       // PostVisibleT 的别名
)

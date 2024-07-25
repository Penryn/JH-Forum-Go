package web

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/alimy/mir/v4"
	"github.com/gin-gonic/gin"
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/model/joint"
	"JH-Forum/internal/servants/base"
	"JH-Forum/pkg/convert"
	"JH-Forum/pkg/xerror"
)

const (
	// 推文可见性
	TweetVisitPublic TweetVisibleType = iota
	TweetVisitPrivate
	TweetVisitFriend
	TweetVisitFollowing
	TweetVisitInvalid
)

// TweetVisibleType 定义了推文可见性类型。
type TweetVisibleType cs.TweetVisibleType

// TweetCommentThumbsReq 定义了动态评论点赞请求结构体。
type TweetCommentThumbsReq struct {
	SimpleInfo `json:"-" binding:"-"`
	TweetId    int64 `json:"tweet_id" binding:"required"`
	CommentId  int64 `json:"comment_id" binding:"required"`
}

// TweetReplyThumbsReq 定义了动态评论回复点赞请求结构体。
type TweetReplyThumbsReq struct {
	SimpleInfo `json:"-" binding:"-"`
	TweetId    int64 `json:"tweet_id" binding:"required"`
	CommentId  int64 `json:"comment_id" binding:"required"`
	ReplyId    int64 `json:"reply_id" binding:"required"`
}

// PostContentItem 定义了动态内容项结构体。
type PostContentItem struct {
	Content string          `json:"content" binding:"required"`
	Type    ms.PostContentT `json:"type" binding:"required"`
	Sort    int64           `json:"sort" binding:"required"`
}

// CreateTweetReq 定义了创建动态请求结构体。
type CreateTweetReq struct {
	BaseInfo   `json:"-" binding:"-"`
	Contents   []*PostContentItem `json:"contents" binding:"required"`
	Tags       []string           `json:"tags" binding:"required"`
	Users      []string           `json:"users" binding:"required"`
	Visibility TweetVisibleType   `json:"visibility"`
	ClientIP   string             `json:"-" binding:"-"`
}

// CreateTweetResp 定义了创建动态响应结构体。
type CreateTweetResp ms.PostFormated

// DeleteTweetReq 定义了删除动态请求结构体。
type DeleteTweetReq struct {
	BaseInfo `json:"-" binding:"-"`
	ID       int64 `json:"id" binding:"required"`
}

// UpateTweetReq 定义了更新动态请求结构体。
type UpdateTweetReq struct {
	BaseInfo   `json:"-" binding:"-"`
	ID         int64              `json:"id" binding:"required"`
	Contents   []*PostContentItem `json:"contents" binding:"required"`
	Tags       []string           `json:"tags" binding:"required"`
	Users      []string           `json:"users" binding:"required"`
	Visibility TweetVisibleType   `json:"visibility"`
	ClientIP   string             `json:"-" binding:"-"`
}

// UpdateTweetResp 定义了更新动态响应结构体。
type UpdateTweetResp ms.PostFormated

// StarTweetReq 定义了动态点赞请求结构体。
type StarTweetReq struct {
	SimpleInfo `json:"-" binding:"-"`
	ID         int64 `json:"id" binding:"required"`
}

// StarTweetResp 定义了动态点赞响应结构体。
type StarTweetResp struct {
	Status bool `json:"status"`
}

// CollectionTweetReq 定义了动态收藏请求结构体。
type CollectionTweetReq struct {
	SimpleInfo `json:"-" binding:"-"`
	ID         int64 `json:"id" binding:"required"`
}

// CollectionTweetResp 定义了动态收藏响应结构体。
type CollectionTweetResp struct {
	Status bool `json:"status"`
}

// LockTweetReq 定义了锁定动态请求结构体。
type LockTweetReq struct {
	BaseInfo `json:"-" binding:"-"`
	ID       int64 `json:"id" binding:"required"`
}

// LockTweetResp 定义了锁定动态响应结构体。
type LockTweetResp struct {
	LockStatus int `json:"lock_status"`
}

// StickTweetReq 定义了置顶动态请求结构体。
type StickTweetReq struct {
	BaseInfo `json:"-" binding:"-"`
	ID       int64 `json:"id" binding:"required"`
}

// HighlightTweetReq 定义了高亮动态请求结构体。
type HighlightTweetReq struct {
	BaseInfo `json:"-" binding:"-"`
	ID       int64 `json:"id" binding:"required"`
}

// StickTweetResp 定义了置顶动态响应结构体。
type StickTweetResp struct {
	StickStatus int `json:"top_status"`
}

// HighlightTweetResp 定义了高亮动态响应结构体。
type HighlightTweetResp struct {
	HighlightStatus int `json:"highlight_status"`
}

// VisibleTweetReq 定义了动态可见性修改请求结构体。
type VisibleTweetReq struct {
	BaseInfo   `json:"-" binding:"-"`
	ID         int64            `json:"id"`
	Visibility TweetVisibleType `json:"visibility"`
}

// VisibleTweetResp 定义了动态可见性修改响应结构体。
type VisibleTweetResp struct {
	Visibility TweetVisibleType `json:"visibility"`
}

// CreateCommentReq 定义了创建评论请求结构体。
type CreateCommentReq struct {
	SimpleInfo `json:"-" binding:"-"`
	PostID     int64              `json:"post_id" binding:"required"`
	Contents   []*PostContentItem `json:"contents" binding:"required"`
	Users      []string           `json:"users" binding:"required"`
	ClientIP   string             `json:"-" binding:"-"`
}

// CreateCommentResp 定义了创建评论响应结构体。
type CreateCommentResp ms.Comment

// CreateCommentReplyReq 定义了创建评论回复请求结构体。
type CreateCommentReplyReq struct {
	SimpleInfo `json:"-" binding:"-"`
	CommentID  int64  `json:"comment_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	AtUserID   int64  `json:"at_user_id"`
	ClientIP   string `json:"-" binding:"-"`
}

// CreateCommentReplyResp 定义了创建评论回复响应结构体。
type CreateCommentReplyResp ms.CommentReply

// DeleteCommentReq 定义了删除评论请求结构体。
type DeleteCommentReq struct {
	BaseInfo `json:"-" binding:"-"`
	ID       int64 `json:"id" binding:"required"`
}

// HighlightCommentReq 定义了高亮评论请求结构体。
type HighlightCommentReq struct {
	SimpleInfo `json:"-" binding:"-"`
	CommentId  int64 `json:"id" binding:"required"`
}

// HighlightCommentResp 定义了高亮评论响应结构体。
type HighlightCommentResp struct {
	HighlightStatus int8 `json:"highlight_status"`
}

// DeleteCommentReplyReq 定义了删除评论回复请求结构体。
type DeleteCommentReplyReq struct {
	BaseInfo `json:"-" binding:"-"`
	ID       int64 `json:"id" binding:"required"`
}

// UploadAttachmentReq 定义了上传附件请求结构体。
type UploadAttachmentReq struct {
	SimpleInfo  `json:"-" binding:"-"`
	UploadType  string
	ContentType string
	File        multipart.File
	FileSize    int64
	FileExt     string
}

// UploadAttachmentResp 定义了上传附件响应结构体。
type UploadAttachmentResp struct {
	UserID    int64             `json:"user_id"`
	FileSize  int64             `json:"file_size"`
	ImgWidth  int               `json:"img_width"`
	ImgHeight int               `json:"img_height"`
	Type      ms.AttachmentType `json:"type"`
	Content   string            `json:"content"`
}

// DownloadAttachmentPrecheckReq 定义了下载附件预检请求结构体。
type DownloadAttachmentPrecheckReq struct {
	BaseInfo  `form:"-" binding:"-"`
	ContentID int64 `form:"id"`
}

// DownloadAttachmentPrecheckResp 定义了下载附件预检响应结构体。
type DownloadAttachmentPrecheckResp struct {
	Paid bool `json:"paid"`
}

// DownloadAttachmentReq 定义了下载附件请求结构体。
type DownloadAttachmentReq struct {
	BaseInfo  `form:"-" binding:"-"`
	ContentID int64 `form:"id"`
}

// DownloadAttachmentResp 定义了下载附件响应结构体。
type DownloadAttachmentResp struct {
	SignedURL string `json:"signed_url"`
}

// StickTopicReq 定义了置顶主题请求结构体。
type StickTopicReq struct {
	SimpleInfo `json:"-" binding:"-"`
	TopicId    int64 `json:"topic_id" binding:"required"`
}

// StickTopicResp 定义了置顶主题响应结构体。
type StickTopicResp struct {
	StickStatus int8 `json:"top_status"`
}

// FollowTopicReq 定义了关注主题请求结构体。
type FollowTopicReq struct {
	SimpleInfo `json:"-" binding:"-"`
	TopicId    int64 `json:"topic_id" binding:"required"`
}

// UnfollowTopicReq 定义了取消关注主题请求结构体。
type UnfollowTopicReq struct {
	SimpleInfo `json:"-" binding:"-"`
	TopicId    int64 `json:"topic_id" binding:"required"`
}

// Check 检查 PostContentItem 属性。
func (p *PostContentItem) Check(acs core.AttachmentCheckService) error {
	// 检查附件是否是本站资源
	if p.Type == ms.ContentTypeImage || p.Type == ms.ContentTypeVideo || p.Type == ms.ContentTypeAttachment {
		if err := acs.CheckAttachment(p.Content); err != nil {
			return err
		}
	}
	// 检查链接是否合法
	if p.Type == ms.ContentTypeLink {
		if !strings.HasPrefix(p.Content, "http://") && !strings.HasPrefix(p.Content, "https://") {
			return fmt.Errorf("链接不合法")
		}
	}
	return nil
}

// Bind 实现了 UploadAttachmentReq 的绑定方法。
func (r *UploadAttachmentReq) Bind(c *gin.Context) (xerr mir.Error) {
	userId, exist := base.UserIdFrom(c)
	if !exist {
		return xerror.UnauthorizedAuthNotExist
	}

	uploadType := c.Request.FormValue("type")
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		return ErrFileUploadFailed
	}
	defer func() {
		if xerr != nil {
			file.Close()
		}
	}()

	if err := fileCheck(uploadType, fileHeader.Size); err != nil {
		return err
	}
	contentType := fileHeader.Header.Get("Content-Type")
	fileExt, xerr := getFileExt(contentType)
	if xerr != nil {
		return xerr
	}
	r.SimpleInfo = SimpleInfo{
		Uid: userId,
	}
	r.UploadType, r.ContentType = uploadType, contentType
	r.File, r.FileSize, r.FileExt = file, fileHeader.Size, fileExt
	return nil
}

// Bind 实现了 DownloadAttachmentPrecheckReq 的绑定方法。
func (r *DownloadAttachmentPrecheckReq) Bind(c *gin.Context) mir.Error {
	user, exist := base.UserFrom(c)
	if !exist {
		return xerror.UnauthorizedAuthNotExist
	}
	r.BaseInfo = BaseInfo{
		User: user,
	}
	r.ContentID = convert.StrTo(c.Query("id")).MustInt64()
	return nil
}

// Bind 实现了 DownloadAttachmentReq 的绑定方法。
func (r *DownloadAttachmentReq) Bind(c *gin.Context) mir.Error {
	user, exist := base.UserFrom(c)
	if !exist {
		return xerror.UnauthorizedAuthNotExist
	}
	r.BaseInfo = BaseInfo{
		User: user,
	}
	r.ContentID = convert.StrTo(c.Query("id")).MustInt64()
	return nil
}

// Bind 实现了 CreateTweetReq 的绑定方法。
func (r *CreateTweetReq) Bind(c *gin.Context) mir.Error {
	r.ClientIP = c.ClientIP()
	return bindAny(c, r)
}

// Bind 实现了 CreateCommentReplyReq 的绑定方法。
func (r *CreateCommentReplyReq) Bind(c *gin.Context) mir.Error {
	r.ClientIP = c.ClientIP()
	return bindAny(c, r)
}

// Bind 实现了 CreateCommentReq 的绑定方法。
func (r *CreateCommentReq) Bind(c *gin.Context) mir.Error {
	r.ClientIP = c.ClientIP()
	return bindAny(c, r)
}

// Render 实现了 CreateTweetResp 的渲染方法。
func (r *CreateTweetResp) Render(c *gin.Context) {
	c.JSON(http.StatusOK, &joint.JsonResp{
		Code: 0,
		Msg:  "success",
		Data: r,
	})
	// 设置审核元信息，用于接下来的审核逻辑
	c.Set(AuditHookCtxKey, &AuditMetaInfo{
		Style: AuditStyleUserTweet,
		Id:    r.ID,
	})
}

// ToVisibleValue 将 TweetVisibleType 转换为内部值。
func (t TweetVisibleType) ToVisibleValue() (res cs.TweetVisibleType) {
	switch t {
	case TweetVisitPublic:
		res = cs.TweetVisitPublic
	case TweetVisitPrivate:
		res = cs.TweetVisitPrivate
	case TweetVisitFriend:
		res = cs.TweetVisitFriend
	case TweetVisitFollowing:
		res = cs.TweetVisitFollowing
	default:
		// 默认私密
		res = cs.TweetVisitPrivate
	}
	return
}

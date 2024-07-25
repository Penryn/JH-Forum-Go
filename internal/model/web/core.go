package web

import (
	"github.com/alimy/mir/v4"
	"github.com/gin-gonic/gin"
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/model/joint"
	"JH-Forum/internal/servants/base"
	"JH-Forum/pkg/convert"
	"JH-Forum/pkg/xerror"
)

type MessageStyle = cs.MessageStyle

// ChangeAvatarReq 定义了修改用户头像请求结构体。
type ChangeAvatarReq struct {
	BaseInfo `json:"-" binding:"-"`
	Avatar   string `json:"avatar" form:"avatar" binding:"required"`
}

// SyncSearchIndexReq 定义了同步搜索索引请求结构体。
type SyncSearchIndexReq struct {
	BaseInfo `json:"-" binding:"-"`
}

// UserInfoReq 定义了获取用户信息请求结构体。
type UserInfoReq struct {
	BaseInfo `json:"-" binding:"-"`
	Username string `json:"username" form:"username" binding:"required"`
}

// UserInfoResp 定义了获取用户信息响应结构体。
type UserInfoResp struct {
	Id          int64  `json:"id"`
	Nickname    string `json:"nickname"`
	Username    string `json:"username"`
	Status      int    `json:"status"`
	Avatar      string `json:"avatar"`
	Balance     int64  `json:"balance"`
	Phone       string `json:"phone"`
	IsAdmin     bool   `json:"is_admin"`
	CreatedOn   int64  `json:"created_on"`
	Follows     int64  `json:"follows"`
	Followings  int64  `json:"followings"`
	TweetsCount int    `json:"tweets_count"`
}

// GetMessagesReq 定义了获取消息请求结构体。
type GetMessagesReq struct {
	SimpleInfo `json:"-" binding:"-"`
	joint.BasePageInfo
	Style MessageStyle `form:"style" binding:"required"`
}

// GetMessagesResp 定义了获取消息响应结构体。
type GetMessagesResp struct {
	joint.CachePageResp
}

// ReadMessageReq 定义了阅读消息请求结构体。
type ReadMessageReq struct {
	SimpleInfo `json:"-" binding:"-"`
	ID         int64 `json:"id" binding:"required"`
}

// ReadAllMessageReq 定义了阅读所有消息请求结构体。
type ReadAllMessageReq struct {
	SimpleInfo `json:"-" binding:"-"`
}

// SendWhisperReq 定义了发送私信请求结构体。
type SendWhisperReq struct {
	SimpleInfo `json:"-" binding:"-"`
	UserID     int64  `json:"user_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// GetCollectionsReq 定义了获取收藏请求结构体。
type GetCollectionsReq BasePageReq

// GetCollectionsResp 定义了获取收藏响应结构体。
type GetCollectionsResp base.PageResp

// GetStarsReq 定义了获取点赞请求结构体。
type GetStarsReq BasePageReq

// GetStarsResp 定义了获取点赞响应结构体。
type GetStarsResp base.PageResp

// UserPhoneBindReq 定义了用户手机绑定请求结构体。
type UserPhoneBindReq struct {
	BaseInfo `json:"-" binding:"-"`
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Captcha  string `json:"captcha" form:"captcha" binding:"required"`
}

// ChangePasswordReq 定义了修改密码请求结构体。
type ChangePasswordReq struct {
	BaseInfo    `json:"-" binding:"-"`
	Password    string `json:"password" form:"password" binding:"required"`
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
}

// ChangeNicknameReq 定义了修改昵称请求结构体。
type ChangeNicknameReq struct {
	BaseInfo `json:"-" binding:"-"`
	Nickname string `json:"nickname" form:"nickname" binding:"required"`
}

// SuggestUsersReq 定义了推荐用户请求结构体。
type SuggestUsersReq struct {
	Keyword string
}

// SuggestUsersResp 定义了推荐用户响应结构体。
type SuggestUsersResp struct {
	Suggests []string `json:"suggest"`
}

// SuggestTagsReq 定义了推荐标签请求结构体。
type SuggestTagsReq struct {
	Keyword string
}

// SuggestTagsResp 定义了推荐标签响应结构体。
type SuggestTagsResp struct {
	Suggests []string `json:"suggest"`
}

// TweetStarStatusReq 定义了获取动态点赞状态请求结构体。
type TweetStarStatusReq struct {
	SimpleInfo `json:"-" binding:"-"`
	TweetId    int64 `form:"id"`
}

// TweetStarStatusResp 定义了获取动态点赞状态响应结构体。
type TweetStarStatusResp struct {
	Status bool `json:"status"`
}

// TweetCollectionStatusReq 定义了获取动态收藏状态请求结构体。
type TweetCollectionStatusReq struct {
	SimpleInfo `json:"-" binding:"-"`
	TweetId    int64 `form:"id"`
}

// TweetCollectionStatusResp 定义了获取动态收藏状态响应结构体。
type TweetCollectionStatusResp struct {
	Status bool `json:"status"`
}

// Bind 实现了 UserInfoReq 结构体的 Bind 方法，用于绑定请求参数到结构体字段。
func (r *UserInfoReq) Bind(c *gin.Context) mir.Error {
	username, exist := base.UserNameFrom(c)
	if !exist {
		return xerror.UnauthorizedAuthNotExist
	}
	r.Username = username
	return nil
}

// Bind 实现了 GetCollectionsReq 结构体的 Bind 方法，用于绑定请求参数到结构体字段。
func (r *GetCollectionsReq) Bind(c *gin.Context) mir.Error {
	return (*BasePageReq)(r).Bind(c)
}

// Bind 实现了 GetStarsReq 结构体的 Bind 方法，用于绑定请求参数到结构体字段。
func (r *GetStarsReq) Bind(c *gin.Context) mir.Error {
	return (*BasePageReq)(r).Bind(c)
}

// Bind 实现了 SuggestTagsReq 结构体的 Bind 方法，用于绑定请求参数到结构体字段。
func (r *SuggestTagsReq) Bind(c *gin.Context) mir.Error {
	r.Keyword = c.Query("k")
	return nil
}

// Bind 实现了 SuggestUsersReq 结构体的 Bind 方法，用于绑定请求参数到结构体字段。
func (r *SuggestUsersReq) Bind(c *gin.Context) mir.Error {
	r.Keyword = c.Query("k")
	return nil
}

// Bind 实现了 TweetCollectionStatusReq 结构体的 Bind 方法，用于绑定请求参数到结构体字段。
func (r *TweetCollectionStatusReq) Bind(c *gin.Context) mir.Error {
	userId, exist := base.UserIdFrom(c)
	if !exist {
		return xerror.UnauthorizedAuthNotExist
	}
	r.SimpleInfo = SimpleInfo{
		Uid: userId,
	}
	r.TweetId = convert.StrTo(c.Query("id")).MustInt64()
	return nil
}

// Bind 实现了 TweetStarStatusReq 结构体的 Bind 方法，用于绑定请求参数到结构体字段。
func (r *TweetStarStatusReq) Bind(c *gin.Context) mir.Error {
	UserId, exist := base.UserIdFrom(c)
	if !exist {
		return xerror.UnauthorizedAuthNotExist
	}
	r.SimpleInfo = SimpleInfo{
		Uid: UserId,
	}
	r.TweetId = convert.StrTo(c.Query("id")).MustInt64()
	return nil
}

package web

import (
	"github.com/alimy/mir/v4"
	"github.com/gin-gonic/gin"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/model/joint"
	"JH-Forum/internal/servants/base"
	"JH-Forum/pkg/app"
)

const (
	TagTypeHot       = cs.TagTypeHot
	TagTypeNew       = cs.TagTypeNew
	TagTypeFollow    = cs.TagTypeFollow
	TagTypeHotExtral = cs.TagTypeHotExtral
)

const (
	UserPostsStylePost      = "post"
	UserPostsStyleComment   = "comment"
	UserPostsStyleHighlight = "highlight"
	UserPostsStyleMedia     = "media"
	UserPostsStyleStar      = "star"

	StyleTweetsNewest    = "newest"
	StyleTweetsHots      = "hots"
	StyleTweetsFollowing = "following"
)

// TagType 定义了标签类型。
type TagType = cs.TagType

// CommentStyleType 定义了评论样式类型。
type CommentStyleType string

// TweetCommentsReq 定义了获取动态评论的请求结构体。
type TweetCommentsReq struct {
	SimpleInfo `form:"-" binding:"-"`
	TweetId    int64            `form:"id" binding:"required"`
	Style      CommentStyleType `form:"style"`
	Page       int              `form:"-" binding:"-"`
	PageSize   int              `form:"-" binding:"-"`
}

// TweetCommentsResp 定义了获取动态评论的响应结构体。
type TweetCommentsResp struct {
	joint.CachePageResp
}

// TimelineReq 定义了获取时间线的请求结构体。
type TimelineReq struct {
	BaseInfo   `form:"-" binding:"-"`
	Query      string              `form:"query"`
	Visibility []core.PostVisibleT `form:"query"`
	Type       string              `form:"type"`
	Style      string              `form:"style"`
	Page       int                 `form:"-" binding:"-"`
	PageSize   int                 `form:"-" binding:"-"`
}

// TimelineResp 定义了获取时间线的响应结构体。
type TimelineResp struct {
	joint.CachePageResp
}

// GetUserTweetsReq 定义了获取用户动态的请求结构体。
type GetUserTweetsReq struct {
	BaseInfo `form:"-" binding:"-"`
	Username string `form:"username" binding:"required"`
	Style    string `form:"style"`
	Page     int    `form:"-" binding:"-"`
	PageSize int    `form:"-" binding:"-"`
}

// GetUserTweetsResp 定义了获取用户动态的响应结构体。
type GetUserTweetsResp struct {
	joint.CachePageResp
}

// GetUserProfileReq 定义了获取用户个人资料的请求结构体。
type GetUserProfileReq struct {
	BaseInfo `form:"-" binding:"-"`
	Username string `form:"username" binding:"required"`
}

// GetUserProfileResp 定义了获取用户个人资料的响应结构体。
type GetUserProfileResp struct {
	ID          int64  `json:"id"`
	Nickname    string `json:"nickname"`
	Username    string `json:"username"`
	Status      int    `json:"status"`
	Avatar      string `json:"avatar"`
	IsAdmin     bool   `json:"is_admin"`
	IsFriend    bool   `json:"is_friend"`
	IsFollowing bool   `json:"is_following"`
	CreatedOn   int64  `json:"created_on"`
	Follows     int64  `json:"follows"`
	Followings  int64  `json:"followings"`
	TweetsCount int    `json:"tweets_count"`
}

// TopicListReq 定义了获取主题列表的请求结构体。
type TopicListReq struct {
	SimpleInfo `form:"-" binding:"-"`
	Type       TagType `json:"type" form:"type" binding:"required"`
	Num        int     `json:"num" form:"num" binding:"required"`
	ExtralNum  int     `json:"extral_num" form:"extral_num"`
}

// TopicListResp 定义了获取主题列表的响应结构体。
type TopicListResp struct {
	Topics       cs.TagList `json:"topics"`
	ExtralTopics cs.TagList `json:"extral_topics,omitempty"`
}

// TweetDetailReq 定义了获取动态详情的请求结构体。
type TweetDetailReq struct {
	SimpleInfo `form:"-" binding:"-"`
	TweetId    int64 `form:"id"`
}

// TweetDetailResp 定义了获取动态详情的响应结构体。
type TweetDetailResp ms.PostFormated

// SetPageInfo 设置分页信息。
func (r *GetUserTweetsReq) SetPageInfo(page int, pageSize int) {
	r.Page, r.PageSize = page, pageSize
}

// SetPageInfo 设置分页信息。
func (r *TweetCommentsReq) SetPageInfo(page int, pageSize int) {
	r.Page, r.PageSize = page, pageSize
}

// Bind 实现了 TimelineReq 的绑定方法。
func (r *TimelineReq) Bind(c *gin.Context) mir.Error {
	user, _ := base.UserFrom(c)
	r.BaseInfo = BaseInfo{
		User: user,
	}
	r.Page, r.PageSize = app.GetPageInfo(c)
	r.Query, r.Type, r.Style = c.Query("query"), "search", c.Query("style")
	return nil
}

// ToInnerValue 将 CommentStyleType 转换为内部值。
func (s CommentStyleType) ToInnerValue() (res cs.StyleCommentType) {
	switch s {
	case "hots":
		res = cs.StyleCommentHots
	case "newest":
		res = cs.StyleCommentNewest
	case "default":
		fallthrough
	default:
		res = cs.StyleCommentDefault
	}
	return
}

// String 返回 CommentStyleType 的字符串表示。
func (s CommentStyleType) String() (res string) {
	switch s {
	case "default":
		res = conf.InfixCommentDefault
	case "hots":
		res = conf.InfixCommentHots
	case "newest":
		res = conf.InfixCommentNewest
	default:
		res = "_"
	}
	return
}

package web

import (
	"github.com/alimy/mir/v4"
	"github.com/gin-gonic/gin"
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/servants/base"
	"JH-Forum/pkg/app"
	"JH-Forum/pkg/xerror"
)

var (
	bindAny = base.NewBindAnyFn()
)

// BaseInfo 包含用户信息的基础结构体。
type BaseInfo struct {
	User *ms.User
}

// SimpleInfo 简单的用户信息结构体。
type SimpleInfo struct {
	Uid int64
}

// BasePageReq 基础分页请求结构体。
type BasePageReq struct {
	UserId   int64
	Page     int
	PageSize int
}

// SetUser 设置用户信息到BaseInfo结构体。
func (b *BaseInfo) SetUser(user *ms.User) {
	b.User = user
}

// SetUserId 设置用户ID到SimpleInfo结构体。
func (s *SimpleInfo) SetUserId(id int64) {
	s.Uid = id
}

// BasePageReqFrom 从Gin上下文中获取基础分页请求信息。
func BasePageReqFrom(c *gin.Context) (*BasePageReq, mir.Error) {
	uid, ok := base.UserIdFrom(c)
	if !ok {
		return nil, xerror.UnauthorizedTokenError
	}
	page, pageSize := app.GetPageInfo(c)
	return &BasePageReq{
		UserId:   uid,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Bind 从Gin上下文中绑定基础分页请求信息到BasePageReq结构体。
func (r *BasePageReq) Bind(c *gin.Context) mir.Error {
	uid, ok := base.UserIdFrom(c)
	if !ok {
		return xerror.UnauthorizedTokenError
	}
	r.UserId = uid
	r.Page, r.PageSize = app.GetPageInfo(c)
	return nil
}

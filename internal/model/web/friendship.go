package web

import (
	"JH-Forum/internal/servants/base"
)

// RequestingFriendReq 定义了请求添加好友的请求结构体。
type RequestingFriendReq struct {
	BaseInfo  `json:"-" binding:"-"`
	UserId    int64  `json:"user_id" binding:"required"`
	Greetings string `json:"greetings" binding:"required"`
}

// AddFriendReq 定义了添加好友的请求结构体。
type AddFriendReq struct {
	BaseInfo `json:"-" binding:"-"`
	UserId   int64 `json:"user_id" binding:"required"`
}

// RejectFriendReq 定义了拒绝好友请求的请求结构体。
type RejectFriendReq struct {
	BaseInfo `json:"-" binding:"-"`
	UserId   int64 `json:"user_id" binding:"required"`
}

// DeleteFriendReq 定义了删除好友的请求结构体。
type DeleteFriendReq struct {
	BaseInfo `json:"-" binding:"-"`
	UserId   int64 `json:"user_id"`
}

// GetContactsReq 定义了获取联系人列表的请求结构体。
type GetContactsReq struct {
	BaseInfo `form:"-" binding:"-"`
	Page     int `form:"-" binding:"-"`
	PageSize int `form:"-" binding:"-"`
}

// GetContactsResp 定义了获取联系人列表的响应结构体。
type GetContactsResp base.PageResp

// SetPageInfo 设置分页信息。
func (r *GetContactsReq) SetPageInfo(page int, pageSize int) {
	r.Page, r.PageSize = page, pageSize
}

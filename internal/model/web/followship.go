package web

import (
	"JH-Forum/internal/model/joint"
	"JH-Forum/internal/servants/base"
)

// FollowUserReq 定义了关注用户请求结构体。
type FollowUserReq struct {
	BaseInfo `json:"-" binding:"-"`
	UserId   int64 `json:"user_id" binding:"required"`
}

// UnfollowUserReq 定义了取消关注用户请求结构体。
type UnfollowUserReq struct {
	BaseInfo `json:"-" binding:"-"`
	UserId   int64 `json:"user_id" binding:"required"`
}

// ListFollowsReq 定义了列出用户关注列表请求结构体。
type ListFollowsReq struct {
	BaseInfo `json:"-" binding:"-"`
	joint.BasePageInfo
	Username string `form:"username" binding:"required"`
}

// ListFollowsResp 定义了列出用户关注列表响应结构体。
type ListFollowsResp base.PageResp

// ListFollowingsReq 定义了列出用户粉丝列表请求结构体。
type ListFollowingsReq struct {
	BaseInfo `form:"-" binding:"-"`
	joint.BasePageInfo
	Username string `form:"username" binding:"required"`
}

// ListFollowingsResp 定义了列出用户粉丝列表响应结构体。
type ListFollowingsResp base.PageResp

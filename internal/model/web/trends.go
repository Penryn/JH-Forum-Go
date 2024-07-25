package web

import (
	"JH-Forum/internal/model/joint"
)

// GetIndexTrendsReq 定义了获取主页趋势请求结构体。
type GetIndexTrendsReq struct {
	SimpleInfo `json:"-" binding:"-"`
	joint.BasePageInfo
}

// GetIndexTrendsResp 定义了获取主页趋势响应结构体。
type GetIndexTrendsResp struct {
	joint.CachePageResp
}

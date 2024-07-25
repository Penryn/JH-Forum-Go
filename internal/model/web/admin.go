package web

// ChangeUserStatusReq 定义了修改用户状态请求结构体。
type ChangeUserStatusReq struct {
	BaseInfo `json:"-" binding:"-"`
	ID       int64 `json:"id" form:"id" binding:"required"`
	Status   int   `json:"status" form:"status" binding:"required,oneof=1 2"`
}

// SiteInfoReq 定义了获取站点信息请求结构体。
type SiteInfoReq struct {
	SimpleInfo `json:"-" binding:"-"`
}

// SiteInfoResp 定义了返回站点信息响应结构体。
type SiteInfoResp struct {
	RegisterUserCount int64 `json:"register_user_count"` // 注册用户数量
	OnlineUserCount   int   `json:"online_user_count"`   // 在线用户数量
	HistoryMaxOnline  int   `json:"history_max_online"`  // 历史最高在线人数
	ServerUpTime      int64 `json:"server_up_time"`      // 服务器运行时间（秒）
}

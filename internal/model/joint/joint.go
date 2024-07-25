// Package joint 提供了一些用于定义模型的常见基础类型或逻辑。

package joint

// BasePageInfo 定义了基础的分页信息结构体。
type BasePageInfo struct {
	Page     int `form:"-" binding:"-"`
	PageSize int `form:"-" binding:"-"`
}

// SetPageInfo 设置分页信息的方法。
func (r *BasePageInfo) SetPageInfo(page int, pageSize int) {
	r.Page, r.PageSize = page, pageSize
}

// JsonResp 定义了用于返回 JSON 响应的结构体。
type JsonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

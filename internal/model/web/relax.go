package web

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"JH-Forum/internal/model/joint"
)

// GetUnreadMsgCountReq 定义了获取未读消息数量请求结构体。
type GetUnreadMsgCountReq struct {
	SimpleInfo `json:"-" binding:"-"`
}

// GetUnreadMsgCountResp 定义了获取未读消息数量响应结构体。
type GetUnreadMsgCountResp struct {
	Count    int64           `json:"count"`
	JsonResp json.RawMessage `json:"-"`
}

// Render 实现了在 Gin 上渲染获取未读消息数量响应的方法。
func (r *GetUnreadMsgCountResp) Render(c *gin.Context) {
	if len(r.JsonResp) != 0 {
		c.JSON(http.StatusOK, r.JsonResp)
	} else {
		c.JSON(http.StatusOK, &joint.JsonResp{
			Code: 0,
			Msg:  "success",
			Data: r,
		})
	}
}

package joint

import (
	stdJson "encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"JH-Forum/pkg/json"
)

// CachePageResp 定义了一个带缓存的页面响应结构体。
type CachePageResp struct {
	Data     *PageResp          // 页面响应数据
	JsonResp stdJson.RawMessage // JSON 响应数据
}

// Render 渲染方法将缓存的响应数据返回给客户端。
func (r *CachePageResp) Render(c *gin.Context) {
	if len(r.JsonResp) != 0 {
		c.JSON(http.StatusOK, r.JsonResp)
	} else {
		c.JSON(http.StatusOK, &JsonResp{
			Code: 0,
			Msg:  "success",
			Data: r.Data,
		})
	}
}

// RespMarshal 将数据序列化为 JSON 格式的字节数组。
func RespMarshal(data interface{}) (stdJson.RawMessage, error) {
	return json.Marshal(data)
}

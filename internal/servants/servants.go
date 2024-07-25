package servants

import (
	"github.com/gin-gonic/gin"
	"JH-Forum/internal/servants/web"
)

// RegisterWebServants register all the servants to gin.Engine
func RegisterWebServants(e *gin.Engine) {
	web.RouteWeb(e)
}

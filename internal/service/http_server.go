package service

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 确保httpServer实现了server接口
var (
	_ server = (*httpServer)(nil)
)

// httpServer结构体，包装了gin.Engine和http.Server
type httpServer struct {
	*baseServer

	e      *gin.Engine
	server *http.Server
}

// 启动httpServer
func (s *httpServer) start() error {
	return s.server.ListenAndServe()
}

// 停止httpServer
func (s *httpServer) stop() error {
	return s.server.Shutdown(context.Background())
}

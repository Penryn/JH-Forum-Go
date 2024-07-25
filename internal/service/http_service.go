package service

import (
	"github.com/gin-gonic/gin"
)

// baseHttpService结构体定义，嵌入了baseService
type baseHttpService struct {
	baseService

	server *httpServer
}

// 注册路由到httpServer
func (s *baseHttpService) registerRoute(srv Service, h func(e *gin.Engine)) {
	if h != nil {
		h(s.server.e)
	}
	s.server.addService(srv)
}

// 在服务启动时的操作，默认不做任何操作
func (s *baseHttpService) OnStart() error {
	// 默认不执行任何操作
	return nil
}

// 在服务停止时的操作，默认不做任何操作
func (s *baseHttpService) OnStop() error {
	// 默认不执行任何操作
	return nil
}

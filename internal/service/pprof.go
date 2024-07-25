package service

import (
	"fmt"
	"net/http"

	_ "net/http/pprof" // 确保这行存在
	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"JH-Forum/internal/conf"
)

// 确保pprofService实现了Service接口
var (
	_ Service = (*pprofService)(nil)
)

// pprofService结构体定义，嵌入了baseHttpService
type pprofService struct {
	*baseHttpService
}

// 返回服务的名称
func (s *pprofService) Name() string {
	return "PprofService"
}

// 返回服务的版本号
func (s *pprofService) Version() *semver.Version {
	return semver.MustParse("v0.1.0")
}

// 初始化函数，注册pprof服务
func (s *pprofService) OnInit() error {
	s.registerRoute(s, nil)
	return nil
}

// 返回服务的描述信息
func (s *pprofService) String() string {
	return fmt.Sprintf("listen on %s\n", color.GreenString("http://%s:%s", conf.PprofServerSetting.HttpIp, conf.PprofServerSetting.HttpPort))
}

// 创建pprofService实例的函数
func newPprofService() Service {
	addr := conf.PprofServerSetting.HttpIp + ":" + conf.PprofServerSetting.HttpPort
	// 创建一个HTTP服务器实例，用于托管pprof服务
	server := httpServers.from(addr, func() *httpServer {
		return &httpServer{
			baseServer: newBaseServe(),
			server: &http.Server{
				Addr:    addr,
				Handler: http.DefaultServeMux, // 确保使用默认的 ServeMux
			},
		}
	})
	// 创建并返回pprofService实例
	return &pprofService{
		baseHttpService: &baseHttpService{
			server: server,
		},
	}
}

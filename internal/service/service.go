package service

import (
	"log"

	"github.com/Masterminds/semver/v3"
	"github.com/alimy/tryst/cfg"
	"JH-Forum/pkg/types"
)

// Service 接口定义了一个服务必须实现的方法
type Service interface {
	Name() string             // 返回服务的名称
	Version() *semver.Version // 返回服务的版本
	OnInit() error            // 服务初始化方法
	OnStart() error           // 服务启动方法
	OnStop() error            // 服务停止方法
}

// baseService 是一个实现了 Service 接口的空服务
type baseService types.Empty

func (baseService) Name() string {
	return ""
}

func (baseService) Version() *semver.Version {
	return semver.MustParse("v0.0.1")
}

func (baseService) String() string {
	return ""
}

// MustInitService 初始化所有服务
func MustInitService() []Service {
	ss := newService() // 创建服务列表
	for _, s := range ss {
		if err := s.OnInit(); err != nil {
			log.Fatalf("initial %s service error: %s", s.Name(), err) // 初始化服务时出错
		}
	}
	return ss
}

// newService 创建所有需要的服务
func newService() (ss []Service) {
	// 根据 config.yaml 中的 features 声明添加所有服务
	cfg.In(cfg.Actions{
		"Web": func() {
			ss = append(ss, newWebService())
		},
		"Pprof": func() {
			ss = append(ss, newPprofService())
		},
	})
	return
}

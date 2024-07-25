package service

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"JH-Forum/internal/conf"
	util "JH-Forum/pkg/utils"
	"github.com/sourcegraph/conc"
)

var (
	httpServers = newServerPool[*httpServer]()
)

const (
	actOnStart byte = iota
	actOnStop
	actStart
	actStop
)

// server 接口定义了服务器的基本操作
type server interface {
	start() error
	stop() error
	services() []Service
}

// serverPool 是服务器对象的池子
type serverPool[T server] struct {
	servers map[string]T
}

// baseServer 实现了服务器的基本操作
type baseServer struct {
	ss map[string]Service
}

// from 根据地址从池子中获取服务器对象，如果不存在则新建
func (p *serverPool[T]) from(addr string, newServer func() T) T {
	s, exist := p.servers[addr]
	if exist {
		return s
	}
	s = newServer()
	p.servers[addr] = s
	return s
}

// startServer 启动所有服务器
func (p *serverPool[T]) startServer(wg *conc.WaitGroup, maxSidSize int) {
	for _, srv := range p.servers {
		ss := srv.services()
		if len(ss) == 0 {
			continue
		}
		startSrv := srv.start
		wg.Go(func() {
			for _, s := range ss {
				colorPrint(actOnStart, s.OnStart(), maxSidSize, s)
			}
			colorPrint(actStart, startSrv(), maxSidSize, ss...)
		})
	}
}

// stopServer 停止所有服务器
func (p *serverPool[T]) stopServer(maxSidSize int) {
	for _, srv := range p.servers {
		ss := srv.services()
		if len(ss) < 1 {
			return
		}
		for _, s := range ss {
			colorPrint(actOnStop, s.OnStop(), maxSidSize, s)
		}
		colorPrint(actStop, srv.stop(), maxSidSize, ss...)
	}
}

// allServices 返回池子中所有的服务对象
func (p *serverPool[T]) allServices() (ss []Service) {
	for _, srv := range p.servers {
		ss = append(ss, srv.services()...)
	}
	return
}

// addService 向 baseServer 中添加服务对象
func (s *baseServer) addService(srv Service) {
	if srv != nil {
		sid := srv.Name() + "@" + srv.Version().String()
		s.ss[sid] = srv
	}
}

// services 返回 baseServer 中所有的服务对象
func (s *baseServer) services() (ss []Service) {
	for _, s := range s.ss {
		ss = append(ss, s)
	}
	return
}

// newServerPool 创建一个新的服务器对象池
func newServerPool[T server]() *serverPool[T] {
	return &serverPool[T]{
		servers: make(map[string]T),
	}
}

// newBaseServe 创建一个新的 baseServer 实例
func newBaseServe() *baseServer {
	return &baseServer{
		ss: make(map[string]Service),
	}
}

// checkServices 检查所有服务的状态和最大服务ID长度
func checkServices() (int, int) {
	var ss []Service
	ss = append(ss, httpServers.allServices()...)
	return len(ss), maxSidSize(ss)
}

// maxSidSize 计算服务ID的最大长度
func maxSidSize(ss []Service) int {
	length := 0
	for _, s := range ss {
		size := len(s.Name() + "@" + s.Version().String())
		if size > length {
			length = size
		}
	}
	return length
}

// colorPrint 根据动作打印彩色日志信息
func colorPrint(act byte, err error, l int, ss ...Service) {
	s := ss[0]
	switch act {
	case actOnStart:
		if err == nil {
			fmt.Fprintf(color.Output, "%s [start] - %s", util.SidStr(s.Name(), s.Version(), l), s)
		} else {
			fmt.Fprintf(color.Output, "%s [start] - 运行 OnStart 出错: %s\n", util.SidStr(s.Name(), s.Version(), l), err)
		}
	case actOnStop:
		if err == nil {
			fmt.Fprintf(color.Output, "%s [stop]  - 完成...\n", util.SidStr(s.Name(), s.Version(), l))
		} else {
			fmt.Fprintf(color.Output, "%s [stop]  - 运行 OnStop 出错: %s\n", util.SidStr(s.Name(), s.Version(), l), err)
		}
	case actStart:
		if err != nil {
			for _, s = range ss {
				fmt.Fprintf(color.Output, "%s [start] - 启动服务器发生错误: %s\n", util.SidStr(s.Name(), s.Version(), l), err)
			}
		}
	case actStop:
		if err != nil {
			for _, s = range ss {
				fmt.Fprintf(color.Output, "%s [stop] - 停止服务器发生错误: %s\n", util.SidStr(s.Name(), s.Version(), l), err)
			}
		}
	}
}

// Start 启动所有服务器
func Start(wg *conc.WaitGroup) {
	srvSize, maxSidSize := checkServices()
	if srvSize < 1 {
		return
	}

	// 为服务器引擎进行一些初始化设置
	gin.SetMode(conf.RunMode())

	// 启动服务器
	httpServers.startServer(wg, maxSidSize)
}

// Stop 停止所有服务器
func Stop() {
	srvSize, maxSidSize := checkServices()
	if srvSize < 1 {
		return
	}
	// 停止服务器
	httpServers.stopServer(maxSidSize)
}

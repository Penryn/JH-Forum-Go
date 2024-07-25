package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"

	"JH-Forum/internal"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/service"
	"JH-Forum/pkg/debug"
	"JH-Forum/pkg/utils"
	"JH-Forum/pkg/version"

	"github.com/sourcegraph/conc"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	// 标志变量：是否不使用默认特性
	noDefaultFeatures bool
	// 标志变量：指定的特性列表
	features []string
)

// serve
func ServeRun() {
	// 打印欢迎横幅，显示版本信息
	utils.PrintHelloBanner(version.VersionInfo())

	// 自动设置最大处理器数量
	maxprocs.Set(maxprocs.Logger(log.Printf))

	// 初始化配置
	conf.Initial(features, noDefaultFeatures)
	internal.Initial()
	ss := service.MustInitService()
	if len(ss) < 1 {
		// 如果没有需要启动的服务，打印提示并退出
		fmt.Fprintln(color.Output, "no service need start so just exit")
		return
	}

	// 如果需要，启动 Pyroscope 进行性能监控
	debug.StartPyroscope()

	// 启动服务，并使用并发控制
	wg := conc.NewWaitGroup()
	fmt.Fprintf(color.Output, "\nstarting run service...\n\n")
	service.Start(wg)

	// 优雅地停止服务
	wg.Go(func() {
		quit := make(chan os.Signal, 1)
		// 监听系统信号，进行优雅停止
		// kill(无参数)默认发送系统调用。SIGTERM
		// kill -2 是系统调用。信号情报
		// kill -9 是系统调用。但是不能被捕获，所以不需要添加它
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Fprintf(color.Output, "\nshutting down server...\n\n")
		service.Stop()
	})
	// 等待所有并发任务完成
	wg.Wait()
}

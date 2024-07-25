package metrics

import (
	"sync" // 导入同步包

	"github.com/alimy/tryst/event" // 导入事件包
	"github.com/alimy/tryst/pool"  // 导入池包
	"JH-Forum/internal/conf"       // 导入配置文件包
	"github.com/sirupsen/logrus"   // 导入日志包
)

var (
	_defaultMetricManager event.EventManager // 默认指标管理器
	_onceInitial          sync.Once          // 一次性初始化锁
)

type Metric = event.Event // 指标类型为事件

type BaseMetric = event.UnimplementedEvent // 基础指标为未实现事件

type MetricManager interface {
	Start()                  // 启动指标管理器
	Stop()                   // 停止指标管理器
	OnMeasure(metric Metric) // 测量指标
}

func StartMetricManager() {
	_defaultMetricManager.Start() // 启动默认指标管理器
}

func StopMetricManager() {
	_defaultMetricManager.Stop() // 停止默认指标管理器
}

// OnMeasure 将指标推送到 goroutine 池中自动处理。
func OnMeasure(metric Metric) {
	_defaultMetricManager.OnEvent(metric) // 推送事件
}

func Initial() {
	_onceInitial.Do(func() { // 仅执行一次初始化操作
		initMetricManager() // 初始化指标管理器
	})
}

func initMetricManager() {
	var opts []pool.Option        // 选项配置数组
	s := conf.EventManagerSetting // 获取事件管理器配置
	if s.MinWorker > 5 {
		opts = append(opts, pool.MinWorkerOpt(s.MinWorker)) // 设置最小工作线程数
	} else {
		opts = append(opts, pool.MinWorkerOpt(5)) // 默认最小工作线程数为5
	}
	if s.MaxEventBuf > 10 {
		opts = append(opts, pool.MaxRequestBufOpt(s.MaxEventBuf)) // 设置最大事件缓冲区
	} else {
		opts = append(opts, pool.MaxRequestBufOpt(10)) // 默认最大事件缓冲区为10
	}
	if s.MaxTempEventBuf > 10 {
		opts = append(opts, pool.MaxRequestTempBufOpt(s.MaxTempEventBuf)) // 设置最大临时事件缓冲区
	} else {
		opts = append(opts, pool.MaxRequestTempBufOpt(10)) // 默认最大临时事件缓冲区为10
	}
	opts = append(opts, pool.MaxTickCountOpt(s.MaxTickCount), pool.TickWaitTimeOpt(s.TickWaitTime)) // 设置最大计数和等待时间
	_defaultMetricManager = event.NewEventManager(func(req Metric, err error) {                     // 创建事件管理器
		if err != nil {
			logrus.Errorf("handle event[%s] occurs error: %s", req.Name(), err) // 处理事件出错时记录错误日志
		}
	}, opts...)
}

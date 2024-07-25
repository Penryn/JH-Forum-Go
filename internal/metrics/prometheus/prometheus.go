package prometheus

import (
	"net/http" // 导入HTTP包

	"github.com/prometheus/client_golang/prometheus"            // 导入Prometheus客户端库
	"github.com/prometheus/client_golang/prometheus/collectors" // 导入Prometheus收集器
	"github.com/prometheus/client_golang/prometheus/promhttp"   // 导入Prometheus HTTP处理器
	"github.com/robfig/cron/v3"                                 // 导入Cron调度库
	"JH-Forum/internal/conf"                                    // 导入配置包
	"JH-Forum/internal/core"                                    // 导入核心包
	"JH-Forum/internal/events"                                  // 导入事件包
	"github.com/sirupsen/logrus"                                // 导入日志库
)

// scheduleJobs 定义调度度量指标更新作业
func scheduleJobs(metrics *metrics) {
	spec := conf.JobManagerSetting.UpdateMetricsInterval // 从配置中获取更新间隔规格
	schedule, err := cron.ParseStandard(spec)            // 解析定时任务规格
	if err != nil {
		panic(err) // 解析失败则抛出异常
	}
	events.OnTask(schedule, metrics.onUpdate)                       // 注册任务调度事件
	logrus.Debug("shedule prometheus metrics update jobs complete") // 记录调试信息
}

// NewHandler 创建Prometheus指标处理器
func NewHandler(ds core.DataService, wc core.WebCache) http.Handler {
	// 创建非全局注册器
	registry := prometheus.NewRegistry()
	// 注册Go运行时指标和进程收集器
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	metrics := newMetrics(registry, ds, wc)                                             // 创建度量指标管理器
	scheduleJobs(metrics)                                                               // 调度度量指标更新作业
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{EnableOpenMetrics: true}) // 返回Prometheus HTTP处理器
}

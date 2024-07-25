package metrics

import (
	"github.com/alimy/tryst/event" // 导入事件包
	"github.com/alimy/tryst/pool"  // 导入池包
)

type simpleMetricManager struct {
	mm event.EventManager // 简单指标管理器
}

func (s *simpleMetricManager) Start() {
	s.mm.Start() // 启动指标管理器
}

func (s *simpleMetricManager) Stop() {
	s.mm.Stop() // 停止指标管理器
}

func (s *simpleMetricManager) OnMeasure(metric Metric) {
	s.mm.OnEvent(metric) // 测量指标
}

// NewMetricManager 创建新的指标管理器
func NewMetricManager(fn pool.RespFn[Metric], opts ...pool.Option) MetricManager {
	return &simpleMetricManager{
		mm: event.NewEventManager(fn, opts...), // 使用事件管理器创建新的指标管理器
	}
}

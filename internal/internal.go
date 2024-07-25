package internal

import (
	"JH-Forum/internal/events"
	"JH-Forum/internal/metrics"
)

// Initial 函数初始化系统的各个组件。
func Initial() {

	// 初始化事件管理系统
	events.Initial()

	// 初始化指标管理系统
	metrics.Initial()
}

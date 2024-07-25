package events

import (
	"github.com/alimy/tryst/event" // 导入事件库
	"github.com/alimy/tryst/pool"  // 导入池库
)

// Event 表示事件类型
type Event = event.Event

// EventManager 定义事件管理器接口
type EventManager interface {
	Start()              // 启动事件管理器
	Stop()               // 停止事件管理器
	OnEvent(event Event) // 处理事件
}

// simpleEventManager 实现简单的事件管理器
type simpleEventManager struct {
	em event.EventManager // 原始事件管理器
}

func (s *simpleEventManager) Start() {
	s.em.Start() // 启动事件管理器
}

func (s *simpleEventManager) Stop() {
	s.em.Stop() // 停止事件管理器
}

func (s *simpleEventManager) OnEvent(event Event) {
	s.em.OnEvent(event) // 处理事件
}

// NewEventManager 创建新的事件管理器
func NewEventManager(fn pool.RespFn[Event], opts ...pool.Option) EventManager {
	return &simpleEventManager{
		em: event.NewEventManager(fn, opts...), // 使用给定的响应函数和选项创建事件管理器
	}
}

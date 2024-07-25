package events

import (
	"sync"

	"github.com/alimy/tryst/cfg"  // 导入配置库
	"github.com/alimy/tryst/pool" // 导入池库
	"github.com/robfig/cron/v3"   // 导入cron库
	"JH-Forum/internal/conf"      // 导入配置包
	"github.com/sirupsen/logrus"  // 导入logrus库
)

var (
	_defaultEventManager EventManager                     // 默认的事件管理器
	_defaultJobManager   JobManager   = emptyJobManager{} // 默认的任务管理器，默认为空实现
	_onceInitial         sync.Once                        // 保证初始化只执行一次的同步锁
)

// StartEventManager 启动默认的事件管理器
func StartEventManager() {
	_defaultEventManager.Start()
}

// StopEventManager 停止默认的事件管理器
func StopEventManager() {
	_defaultEventManager.Stop()
}

// OnEvent 将事件推送到协程池中自动处理
func OnEvent(event Event) {
	_defaultEventManager.OnEvent(event)
}

// StartJobManager 启动默认的任务管理器
func StartJobManager() {
	_defaultJobManager.Start()
}

// StopJobManager 停止默认的任务管理器
func StopJobManager() {
	_defaultJobManager.Stop()
}

// NewJob 创建一个新的Job实例
func NewJob(s cron.Schedule, fn JobFn) Job {
	return &simpleJob{
		Schedule: s,
		Job:      fn,
	}
}

// RemoveJob 从未来的运行中移除一个条目
func RemoveJob(id EntryID) {
	_defaultJobManager.Remove(id)
}

// Schedule 将一个Job添加到Cron中，按给定的调度运行
func Schedule(job Job) EntryID {
	return _defaultJobManager.Schedule(job)
}

// OnTask 将一个Job添加到Cron中，按给定的调度运行
func OnTask(s cron.Schedule, fn JobFn) EntryID {
	job := &simpleJob{
		Schedule: s,
		Job:      fn,
	}
	return _defaultJobManager.Schedule(job)
}

// Initial 初始化事件管理器和任务管理器
func Initial() {
	_onceInitial.Do(func() {
		initEventManager() // 初始化事件管理器
		cfg.Not("DisableJobManager", func() {
			initJobManager() // 初始化任务管理器
			logrus.Debugln("initial JobManager")
		})
	})
}

// initJobManager 初始化任务管理器
func initJobManager() {
	_defaultJobManager = NewJobManager() // 创建新的任务管理器实例
	StartJobManager()                    // 启动任务管理器
}

// initEventManager 初始化事件管理器
func initEventManager() {
	var opts []pool.Option
	s := conf.EventManagerSetting // 获取事件管理器的配置
	if s.MinWorker > 5 {
		opts = append(opts, pool.MinWorkerOpt(s.MinWorker))
	} else {
		opts = append(opts, pool.MinWorkerOpt(5))
	}
	if s.MaxEventBuf > 10 {
		opts = append(opts, pool.MaxRequestBufOpt(s.MaxEventBuf))
	} else {
		opts = append(opts, pool.MaxRequestBufOpt(10))
	}
	if s.MaxTempEventBuf > 10 {
		opts = append(opts, pool.MaxRequestTempBufOpt(s.MaxTempEventBuf))
	} else {
		opts = append(opts, pool.MaxRequestTempBufOpt(10))
	}
	opts = append(opts, pool.MaxTickCountOpt(s.MaxTickCount), pool.TickWaitTimeOpt(s.TickWaitTime))
	_defaultEventManager = NewEventManager(func(req Event, err error) {
		if err != nil {
			logrus.Errorf("handle event[%s] occurs error: %s", req.Name(), err)
		}
	}, opts...)
}

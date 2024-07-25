package events

import (
	"github.com/robfig/cron/v3" // 导入cron库
	"JH-Forum/pkg/types"        // 导入types包
)

type (
	EntryID = cron.EntryID // 定义EntryID类型，基于cron库的EntryID
)

// JobFn 定义Job函数类型，实现cron.Job接口
type JobFn func()

func (fn JobFn) Run() {
	fn()
}

// Job 定义Job接口，组合了cron.Schedule和cron.Job接口
type Job interface {
	cron.Schedule
	cron.Job
}

type simpleJob struct {
	cron.Schedule
	cron.Job
}

// JobManager 定义Job管理器接口
type JobManager interface {
	Start()
	Stop()
	Remove(id EntryID)
	Schedule(Job) EntryID
}

// emptyJobManager 是一个空实现的JobManager
type emptyJobManager types.Empty

// simpleJobManager 是一个简单实现的JobManager
type simpleJobManager struct {
	m *cron.Cron // cron调度器实例
}

func (emptyJobManager) Start() {
	// 空实现
}

func (emptyJobManager) Stop() {
	// 空实现
}

func (emptyJobManager) Remove(id EntryID) {
	// 空实现
}

func (emptyJobManager) Schedule(job Job) EntryID {
	return 0
}

func (j *simpleJobManager) Start() {
	j.m.Start()
}

func (j *simpleJobManager) Stop() {
	j.m.Stop()
}

// Remove 从未来的运行中移除一个条目
func (j *simpleJobManager) Remove(id EntryID) {
	j.m.Remove(id)
}

// Schedule 将一个Job添加到Cron中，按给定的调度运行
func (j *simpleJobManager) Schedule(job Job) EntryID {
	return j.m.Schedule(job, job)
}

// NewJobManager 创建一个新的JobManager实例
func NewJobManager() JobManager {
	return &simpleJobManager{
		m: cron.New(), // 创建一个新的cron调度器实例
	}
}

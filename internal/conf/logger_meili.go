package conf

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/sirupsen/logrus"
)

// meiliLogData 是一个包含 map[string]any 类型的切片
type meiliLogData []map[string]any

// meiliLogHook 是一个实现了 logrus.Hook 接口的结构体
type meiliLogHook struct {
	config    meilisearch.ClientConfig // MeiliSearch 客户端配置
	idxName   string                   // 索引名称
	addDocsCh chan *meiliLogData       // 日志数据通道
}

// Fire 方法实现了 logrus.Hook 接口的 Fire 方法，用于处理日志条目
func (h *meiliLogHook) Fire(entry *logrus.Entry) error {
	data := meiliLogData{{
		"id":      entry.Time.Unix(),
		"time":    entry.Time,
		"level":   entry.Level,
		"message": entry.Message,
		"data":    entry.Data,
	}}

	// 尝试发送日志数据到 addDocsCh，如果通道已满则新开 goroutine 加入文档
	select {
	case h.addDocsCh <- &data:
	default:
		go func(index *meilisearch.Index, item meiliLogData) {
			index.AddDocuments(item)
		}(h.index(), data)
	}

	return nil
}

// handleAddDocs 方法用于处理从 addDocsCh 中接收的日志数据并加入到 MeiliSearch 索引中
func (h *meiliLogHook) handleAddDocs() {
	index := h.index()
	for item := range h.addDocsCh {
		index.AddDocuments(item)
	}
}

// Levels 方法实现了 logrus.Hook 接口的 Levels 方法，返回支持的日志级别
func (h *meiliLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// index 方法返回 MeiliSearch 索引对象
func (h *meiliLogHook) index() *meilisearch.Index {
	return meilisearch.NewClient(h.config).Index(h.idxName)
}

// newMeiliLogHook 用于创建一个新的 meiliLogHook 实例
func newMeiliLogHook() *meiliLogHook {
	hook := &meiliLogHook{
		config: meilisearch.ClientConfig{
			Host:   loggerMeiliSetting.Endpoint(),
			APIKey: loggerMeiliSetting.ApiKey,
		},
		idxName: loggerMeiliSetting.Index,
	}

	client := meilisearch.NewClient(hook.config)
	index := client.Index(hook.idxName)
	if _, err := index.FetchInfo(); err != nil {
		client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        hook.idxName,
			PrimaryKey: "id",
		})
		sortableAttributes := []string{
			"time",
		}
		index.UpdateSortableAttributes(&sortableAttributes)
	}

	// 初始化 addDocsCh
	hook.addDocsCh = make(chan *meiliLogData, loggerMeiliSetting.maxLogBuffer())

	// 启动后台日志处理 goroutine
	for minWork := loggerMeiliSetting.minWork(); minWork > 0; minWork-- {
		go hook.handleAddDocs()
	}

	return hook
}

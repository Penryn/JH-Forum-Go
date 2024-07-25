// 包 search 实现了与不同搜索引擎集成的 Tweet 搜索服务。

package search

import (
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"github.com/sirupsen/logrus"
)

// NewMeiliTweetSearchService 创建一个基于 MeiliSearch 的 Tweet 搜索服务实例。
func NewMeiliTweetSearchService(ams core.AuthorizationManageService) (core.TweetSearchService, core.VersionInfo) {
	s := conf.MeiliSetting
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   s.Endpoint(),
		APIKey: s.ApiKey,
	})

	// 如果索引不存在，则创建新索引并设置初始设置
	if _, err := client.Index(s.Index).FetchInfo(); err != nil {
		logrus.Debugf("create meili index because fetch index info error: %v", err)
		if _, err := client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        s.Index,
			PrimaryKey: "id",
		}); err == nil {
			settings := meilisearch.Settings{
				SearchableAttributes: []string{"content", "tags"},
				SortableAttributes:   []string{"is_top", "latest_replied_on"},
				FilterableAttributes: []string{"tags", "visibility", "user_id"},
			}
			if _, err = client.Index(s.Index).UpdateSettings(&settings); err != nil {
				logrus.Errorf("update meili settings error: %s", err)
			}
		} else {
			logrus.Errorf("create meili index error: %s", err)
		}
	}

	// 创建 MeiliTweetSearchServant 实例
	mts := &meiliTweetSearchServant{
		tweetSearchFilter: tweetSearchFilter{
			ams: ams,
		},
		client:        client,
		index:         client.Index(s.Index),
		publicFilter:  fmt.Sprintf("visibility=%d", core.PostVisitPublic),
		privateFilter: fmt.Sprintf("visibility=%d AND user_id=", core.PostVisitPrivate),
		friendFilter:  fmt.Sprintf("visibility=%d", core.PostVisitFriend),
	}
	return mts, mts
}

// NewBridgeTweetSearchService 创建一个桥接不同 Tweet 搜索服务的服务实例。
func NewBridgeTweetSearchService(ts core.TweetSearchService) core.TweetSearchService {
	capacity := conf.TweetSearchSetting.MaxUpdateQPS
	if capacity < 10 {
		capacity = 10
	} else if capacity > 10000 {
		capacity = 10000
	}
	bts := &bridgeTweetSearchServant{
		ts:               ts,
		updateDocsCh:     make(chan *documents, capacity),
		updateDocsTempCh: make(chan *documents, 100),
	}

	numWorker := conf.TweetSearchSetting.MinWorker
	if numWorker < 5 {
		numWorker = 5
	} else if numWorker > 1000 {
		numWorker = 1000
	}
	logrus.Debugf("use %d backend worker to update documents to search engine")
	// 启动文档更新器
	for ; numWorker > 0; numWorker-- {
		go bts.startUpdateDocs()
	}

	return bts
}

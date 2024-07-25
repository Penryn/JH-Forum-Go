// 包 search 实现了与 Tweet 相关的搜索服务，通过桥接模式与具体的 Tweet 搜索服务进行交互和更新。

package search

import (
	"time"

	"JH-Forum/internal/core"
	"JH-Forum/internal/core/ms"
	"github.com/sirupsen/logrus"
)

// documents 定义了索引文档的结构。
type documents struct {
	primaryKey  []string         // 主键
	docItems    []core.TsDocItem // 文档项
	identifiers []string         // 标识符
}

// bridgeTweetSearchServant 实现了 Tweet 搜索服务的桥接。
type bridgeTweetSearchServant struct {
	ts               core.TweetSearchService // 具体的 Tweet 搜索服务
	updateDocsCh     chan *documents         // 更新文档的通道
	updateDocsTempCh chan *documents         // 临时更新文档的通道
}

// IndexName 返回搜索服务的索引名称。
func (s *bridgeTweetSearchServant) IndexName() string {
	return s.ts.IndexName()
}

// AddDocuments 向索引中添加文档。
func (s *bridgeTweetSearchServant) AddDocuments(data []core.TsDocItem, primaryKey ...string) (bool, error) {
	s.updateDocs(&documents{
		primaryKey: primaryKey,
		docItems:   data,
	})
	return true, nil
}

// DeleteDocuments 删除索引中的文档。
func (s *bridgeTweetSearchServant) DeleteDocuments(identifiers []string) error {
	s.updateDocs(&documents{
		identifiers: identifiers,
	})
	return nil
}

// Search 执行搜索操作。
func (s *bridgeTweetSearchServant) Search(user *ms.User, q *core.QueryReq, offset, limit int) (*core.QueryResp, error) {
	return s.ts.Search(user, q, offset, limit)
}

// updateDocs 更新文档。
func (s *bridgeTweetSearchServant) updateDocs(doc *documents) {
	select {
	case s.updateDocsCh <- doc:
		logrus.Debugln("addDocuments send documents by updateDocsCh chan")
	default:
		select {
		case s.updateDocsTempCh <- doc:
			logrus.Debugln("addDocuments send documents by updateDocsTempCh chan")
		default:
			go func() {
				s.handleUpdate(doc)

				// 监视 updateDocsTempCh，继续处理更新（如果需要的话）。
				// 如果一分钟内没有处理任何项目，则取消循环。
				for count := 0; count < 60; count++ {
					select {
					case item := <-s.updateDocsTempCh:
						// 重置计数以继续处理文档更新
						count = 0
						s.handleUpdate(item)
					default:
						// 等待文档项传递以处理
						time.Sleep(1 * time.Second)
					}
				}
			}()
		}
	}
}

// startUpdateDocs 启动文档更新处理。
func (s *bridgeTweetSearchServant) startUpdateDocs() {
	for doc := range s.updateDocsCh {
		s.handleUpdate(doc)
	}
}

// handleUpdate 处理文档更新。
func (s *bridgeTweetSearchServant) handleUpdate(item *documents) {
	if len(item.docItems) > 0 {
		if _, err := s.ts.AddDocuments(item.docItems, item.primaryKey...); err != nil {
			logrus.Errorf("addDocuments 发生错误: %v", err)
		}
	} else if len(item.identifiers) > 0 {
		if err := s.ts.DeleteDocuments(item.identifiers); err != nil {
			logrus.Errorf("deleteDocuments 发生错误: %s", err)
		}
	}
}

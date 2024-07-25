// 该代码定义了推文搜索服务的接口和相关类型，用于实现推文的索引、删除和搜索功能。

package core

import (
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/dao/jinzhu/dbr"
)

// 搜索类型常量
const (
	SearchTypeDefault SearchType = "search"
	SearchTypeTag     SearchType = "tag"
)

// 帖子访问权限常量
const (
	PostVisitPublic    = dbr.PostVisitPublic
	PostVisitPrivate   = dbr.PostVisitPrivate
	PostVisitFriend    = dbr.PostVisitFriend
	PostVisitFollowing = dbr.PostVisitFollowing
)

type (
	// PostVisibleT 可访问类型，可见性: 0私密 10充电可见 20订阅可见 30保留 40保留 50好友可见 60关注可见 70保留 80保留 90公开
	PostVisibleT = dbr.PostVisibleT

	SearchType string

	// QueryReq 查询请求结构体
	QueryReq struct {
		Query      string         // 查询内容
		Visibility []PostVisibleT // 可见性过滤
		Type       SearchType     // 搜索类型
	}

	// QueryResp 查询响应结构体
	QueryResp struct {
		Items []*ms.PostFormated // 查询结果项
		Total int64              // 结果总数
	}

	// TsDocItem 文档项结构体
	TsDocItem struct {
		Post    *ms.Post // 推文
		Content string   // 推文内容
	}
)

// TweetSearchService 推文搜索服务接口
type TweetSearchService interface {
	IndexName() string                                                        // 获取索引名称
	AddDocuments(data []TsDocItem, primaryKey ...string) (bool, error)        // 添加文档
	DeleteDocuments(identifiers []string) error                               // 删除文档
	Search(user *ms.User, q *QueryReq, offset, limit int) (*QueryResp, error) // 搜索文档
}

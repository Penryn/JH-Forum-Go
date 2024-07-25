// 该接口定义了缓存索引服务、Redis缓存和应用缓存的功能。

package core

import (
	"context"

	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/dao/jinzhu/dbr"
)

// 索引动作的常量
const (
	IdxActNop IdxAct = iota + 1
	IdxActCreatePost
	IdxActUpdatePost
	IdxActDeletePost
	IdxActStickPost
	IdxActVisiblePost
)

type IdxAct uint8

// 索引动作结构体
type IndexAction struct {
	Act  IdxAct
	Post *dbr.Post
}

// 索引动作结构体A
type IndexActionA struct {
	Act   IdxAct
	Tweet *cs.TweetInfo
}

// 将索引动作转换为字符串
func (a IdxAct) String() string {
	switch a {
	case IdxActNop:
		return "no operator"
	case IdxActCreatePost:
		return "create post"
	case IdxActUpdatePost:
		return "update post"
	case IdxActDeletePost:
		return "delete post"
	case IdxActStickPost:
		return "stick post"
	case IdxActVisiblePost:
		return "visible post"
	default:
		return "unknow action"
	}
}

// 创建新的索引动作
func NewIndexAction(act IdxAct, post *ms.Post) *IndexAction {
	return &IndexAction{
		Act:  act,
		Post: post,
	}
}

// 创建新的索引动作A
func NewIndexActionA(act IdxAct, tweet *cs.TweetInfo) *IndexActionA {
	return &IndexActionA{
		Act:   act,
		Tweet: tweet,
	}
}

// CacheIndexService 缓存索引服务接口
type CacheIndexService interface {
	// 发送索引动作
	SendAction(act IdxAct, post *dbr.Post)
}

// CacheIndexServantA 缓存索引服务接口A
type CacheIndexServantA interface {
	// 发送索引动作
	SendAction(act IdxAct, tweet *cs.TweetInfo)
}

// RedisCache 基于Redis的内存缓存
type RedisCache interface {
	SetPushToSearchJob(ctx context.Context) error
	DelPushToSearchJob(ctx context.Context) error
	SetImgCaptcha(ctx context.Context, id string, value string) error
	GetImgCaptcha(ctx context.Context, id string) (string, error)
	DelImgCaptcha(ctx context.Context, id string) error
	GetCountLoginErr(ctx context.Context, id int64) (int64, error)
	DelCountLoginErr(ctx context.Context, id int64) error
	IncrCountLoginErr(ctx context.Context, id int64) error
	GetCountWhisper(ctx context.Context, uid int64) (int64, error)
	IncrCountWhisper(ctx context.Context, uid int64) error
	SetRechargeStatus(ctx context.Context, tradeNo string) error
	DelRechargeStatus(ctx context.Context, tradeNo string) error
}

// AppCache 应用缓存接口
type AppCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte, ex int64) error
	SetNx(key string, data []byte, ex int64) error
	Delete(key ...string) error
	DelAny(pattern string) error
	Exist(key string) bool
	Keys(pattern string) ([]string, error)
}

// WebCache 网络缓存接口
type WebCache interface {
	AppCache
	GetUnreadMsgCountResp(uid int64) ([]byte, error)
	PutUnreadMsgCountResp(uid int64, data []byte) error
	DelUnreadMsgCountResp(uid int64) error
	ExistUnreadMsgCountResp(uid int64) bool
	PutHistoryMaxOnline(newScore int) (int, error)
}

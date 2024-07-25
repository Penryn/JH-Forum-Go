// Copyright 2023 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"context"
	"time"

	"github.com/redis/rueidis"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"JH-Forum/pkg/utils"
)

var (
	_webCache core.WebCache = (*webCache)(nil)
	_appCache core.AppCache = (*appCache)(nil)
)

// appCache 实现了 AppCache 接口，用于管理应用程序级别的缓存
type appCache struct {
	cscExpire time.Duration
	c         rueidis.Client
}

type webCache struct {
	core.AppCache
	c               rueidis.Client
	unreadMsgExpire int64
}

// Get 获取指定键的缓存数据
func (s *appCache) Get(key string) ([]byte, error) {
	res, err := rueidis.MGetCache(s.c, context.Background(), s.cscExpire, []string{key})
	if err != nil {
		return nil, err
	}
	message := res[key]
	return message.AsBytes()
}

// Set 设置指定键的缓存数据，可以指定过期时间（秒）
func (s *appCache) Set(key string, data []byte, ex int64) error {
	ctx := context.Background()
	cmd := s.c.B().Set().Key(key).Value(utils.String(data))
	if ex > 0 {
		return s.c.Do(ctx, cmd.ExSeconds(ex).Build()).Error()
	}
	return s.c.Do(ctx, cmd.Build()).Error()
}

// SetNx 设置指定键的缓存数据，仅当键不存在时有效，可以指定过期时间（秒）
func (s *appCache) SetNx(key string, data []byte, ex int64) error {
	ctx := context.Background()
	cmd := s.c.B().Set().Key(key).Value(utils.String(data)).Nx()
	if ex > 0 {
		return s.c.Do(ctx, cmd.ExSeconds(ex).Build()).Error()
	}
	return s.c.Do(ctx, cmd.Build()).Error()
}

// Delete 删除指定键的缓存数据
func (s *appCache) Delete(keys ...string) (err error) {
	if len(keys) != 0 {
		err = s.c.Do(context.Background(), s.c.B().Del().Key(keys...).Build()).Error()
	}
	return
}

// DelAny 根据指定的模式删除匹配的缓存键
func (s *appCache) DelAny(pattern string) (err error) {
	var (
		keys   []string
		cursor uint64
		entry  rueidis.ScanEntry
	)
	ctx := context.Background()
	for {
		cmd := s.c.B().Scan().Cursor(cursor).Match(pattern).Count(50).Build()
		if entry, err = s.c.Do(ctx, cmd).AsScanEntry(); err != nil {
			return
		}
		keys = append(keys, entry.Elements...)
		if entry.Cursor != 0 {
			cursor = entry.Cursor
			continue
		}
		break
	}
	if len(keys) != 0 {
		err = s.c.Do(ctx, s.c.B().Del().Key(keys...).Build()).Error()
	}
	return
}

// Exist 检查指定键是否存在
func (s *appCache) Exist(key string) bool {
	cmd := s.c.B().Exists().Key(key).Build()
	count, _ := s.c.Do(context.Background(), cmd).AsInt64()
	return count > 0
}

// Keys 获取所有匹配指定模式的键
func (s *appCache) Keys(pattern string) (res []string, err error) {
	ctx, cursor := context.Background(), uint64(0)
	for {
		cmd := s.c.B().Scan().Cursor(cursor).Match(pattern).Count(50).Build()
		entry, err := s.c.Do(ctx, cmd).AsScanEntry()
		if err != nil {
			return nil, err
		}
		res = append(res, entry.Elements...)
		if entry.Cursor != 0 {
			cursor = entry.Cursor
			continue
		}
		break
	}
	return
}

// GetUnreadMsgCountResp 获取指定用户未读消息计数的缓存数据
func (s *webCache) GetUnreadMsgCountResp(uid int64) ([]byte, error) {
	key := conf.KeyUnreadMsg.Get(uid)
	return s.Get(key)
}

// PutUnreadMsgCountResp 设置指定用户未读消息计数的缓存数据
func (s *webCache) PutUnreadMsgCountResp(uid int64, data []byte) error {
	return s.Set(conf.KeyUnreadMsg.Get(uid), data, s.unreadMsgExpire)
}

// DelUnreadMsgCountResp 删除指定用户未读消息计数的缓存数据
func (s *webCache) DelUnreadMsgCountResp(uid int64) error {
	return s.Delete(conf.KeyUnreadMsg.Get(uid))
}

// ExistUnreadMsgCountResp 检查指定用户未读消息计数的缓存数据是否存在
func (s *webCache) ExistUnreadMsgCountResp(uid int64) bool {
	return s.Exist(conf.KeyUnreadMsg.Get(uid))
}

// PutHistoryMaxOnline 设置历史最大在线用户数，并返回设置后的值
func (s *webCache) PutHistoryMaxOnline(newScore int) (int, error) {
	ctx := context.Background()
	cmd := s.c.B().Zadd().
		Key(conf.KeySiteStatus).
		Gt().ScoreMember().
		ScoreMember(float64(newScore), conf.KeyHistoryMaxOnline).Build()
	if err := s.c.Do(ctx, cmd).Error(); err != nil {
		return 0, err
	}
	cmd = s.c.B().Zscore().Key(conf.KeySiteStatus).Member(conf.KeyHistoryMaxOnline).Build()
	if score, err := s.c.Do(ctx, cmd).ToFloat64(); err == nil {
		return int(score), nil
	} else {
		return 0, err
	}
}

// newAppCache 创建一个新的 appCache 实例
func newAppCache() *appCache {
	return &appCache{
		cscExpire: conf.CacheSetting.CientSideCacheExpire,
		c:         conf.MustRedisClient(),
	}
}

// newWebCache 创建一个新的 webCache 实例
func newWebCache(ac core.AppCache) *webCache {
	return &webCache{
		AppCache:        ac,
		c:               conf.MustRedisClient(),
		unreadMsgExpire: conf.CacheSetting.UnreadMsgExpire,
	}
}

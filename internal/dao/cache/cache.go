package cache

import (
	"context"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"github.com/sirupsen/logrus"
)

var (
	_onceInit sync.Once
)

// NewRedisCache 创建一个 Redis 缓存实例
func NewRedisCache() core.RedisCache {
	return &redisCache{
		c: conf.MustRedisClient(),
	}
}

// NewBigCacheIndexService 创建一个基于 BigCache 的索引缓存服务
func NewBigCacheIndexService(ips core.IndexPostsService, ams core.AuthorizationManageService) (core.CacheIndexService, core.VersionInfo) {
	s := conf.BigCacheIndexSetting
	c := bigcache.DefaultConfig(s.ExpireInSecond)
	c.Shards = s.MaxIndexPage
	c.HardMaxCacheSize = s.HardMaxCacheSize
	c.Verbose = s.Verbose
	c.MaxEntrySize = 10000
	c.Logger = logrus.StandardLogger()

	bc, err := bigcache.New(context.Background(), c)
	if err != nil {
		logrus.Fatalf("初始化 BigCacheIndex 失败，错误：%v", err)
	}
	cacheIndex := newCacheIndexSrv(ips, ams, &bigCacheTweetsCache{
		bc: bc,
	})
	return cacheIndex, cacheIndex
}

// NewRedisCacheIndexService 创建一个基于 Redis 的索引缓存服务
func NewRedisCacheIndexService(ips core.IndexPostsService, ams core.AuthorizationManageService) (core.CacheIndexService, core.VersionInfo) {
	cacheIndex := newCacheIndexSrv(ips, ams, &redisCacheTweetsCache{
		expireDuration: conf.RedisCacheIndexSetting.ExpireInSecond,
		expireInSecond: int64(conf.RedisCacheIndexSetting.ExpireInSecond / time.Second),
		c:              conf.MustRedisClient(),
	})
	return cacheIndex, cacheIndex
}

// NewWebCache 返回一个 Web 缓存实例
func NewWebCache() core.WebCache {
	lazyInitial()
	return _webCache
}

// NewAppCache 返回一个 App 缓存实例
func NewAppCache() core.AppCache {
	lazyInitial()
	return _appCache
}

// NewSimpleCacheIndexService 创建一个简单的索引缓存服务
func NewSimpleCacheIndexService(indexPosts core.IndexPostsService) (core.CacheIndexService, core.VersionInfo) {
	s := conf.SimpleCacheIndexSetting
	cacheIndex := &simpleCacheIndexServant{
		ips:             indexPosts,
		maxIndexSize:    s.MaxIndexSize,
		indexPosts:      nil,
		checkTick:       time.NewTicker(s.CheckTickDuration), // 每分钟检查是否需要更新索引
		expireIndexTick: time.NewTicker(time.Second),
	}

	// 每 ExpireTickDuration 秒强制过期索引
	if s.ExpireTickDuration != 0 {
		cacheIndex.expireIndexTick.Reset(s.CheckTickDuration)
	} else {
		cacheIndex.expireIndexTick.Stop()
	}

	// indexActionCh 的容量可以通过 conf.yaml 进行配置，需在 [10, 10000] 范围内，或重新编译源码以调整最小/最大容量
	capacity := conf.CacheIndexSetting.MaxUpdateQPS
	if capacity < 10 {
		capacity = 10
	} else if capacity > 10000 {
		capacity = 10000
	}
	cacheIndex.indexActionCh = make(chan core.IdxAct, capacity)

	// 启动索引更新
	cacheIndex.atomicIndex.Store(cacheIndex.indexPosts)
	go cacheIndex.startIndexPosts()

	return cacheIndex, cacheIndex
}

// NewNoneCacheIndexService 创建一个不使用缓存的索引缓存服务
func NewNoneCacheIndexService(indexPosts core.IndexPostsService) (core.CacheIndexService, core.VersionInfo) {
	obj := &noneCacheIndexServant{
		ips: indexPosts,
	}
	return obj, obj
}

// lazyInitial 实现延迟初始化逻辑
func lazyInitial() {
	_onceInit.Do(func() {
		_appCache = newAppCache()
		_webCache = newWebCache(_appCache)
	})
}

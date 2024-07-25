package cache

import (
	"github.com/Masterminds/semver/v3"
	"github.com/allegro/bigcache/v3"
)

var (
	_ tweetsCache = (*bigCacheTweetsCache)(nil)
)

// bigCacheTweetsCache 实现了 tweetsCache 接口
type bigCacheTweetsCache struct {
	bc *bigcache.BigCache
}

// getTweetsBytes 从 BigCache 获取推文数据
func (s *bigCacheTweetsCache) getTweetsBytes(key string) ([]byte, error) {
	return s.bc.Get(key)
}

// setTweetsBytes 将推文数据写入 BigCache
func (s *bigCacheTweetsCache) setTweetsBytes(key string, bs []byte) error {
	return s.bc.Set(key, bs)
}

// delTweets 删除指定键的推文数据
func (s *bigCacheTweetsCache) delTweets(keys []string) error {
	for _, k := range keys {
		s.bc.Delete(k)
	}
	return nil
}

// allKeys 获取所有缓存键
func (s *bigCacheTweetsCache) allKeys() ([]string, error) {
	var keys []string
	for it := s.bc.Iterator(); it.SetNext(); {
		entry, err := it.Value()
		if err != nil {
			return nil, err
		}
		keys = append(keys, entry.Key())
	}
	return keys, nil
}

// Name 返回缓存名称
func (s *bigCacheTweetsCache) Name() string {
	return "BigCacheIndex"
}

// Version 返回缓存版本信息
func (s *bigCacheTweetsCache) Version() *semver.Version {
	return semver.MustParse("v0.2.0")
}

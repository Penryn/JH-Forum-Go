// Copyright 2023 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/ms"
	"JH-Forum/pkg/types"
	"github.com/sirupsen/logrus"
)

// _cacheIndexKey 是缓存中索引的键前缀
const (
	_cacheIndexKey = "paopao_index"
)

var (
	_ core.CacheIndexService = (*cacheIndexSrv)(nil)
	_ core.VersionInfo       = (*cacheIndexSrv)(nil)
)

type postsEntry struct {
	key    string
	tweets *ms.IndexTweetList
}

type tweetsCache interface {
	core.VersionInfo
	getTweetsBytes(key string) ([]byte, error)
	setTweetsBytes(key string, bs []byte) error
	delTweets(keys []string) error
	allKeys() ([]string, error)
}

// cacheIndexSrv 实现了 core.CacheIndexService 和 core.VersionInfo 接口
type cacheIndexSrv struct {
	ips core.IndexPostsService          // 索引推文服务接口
	ams core.AuthorizationManageService // 授权管理服务接口

	name               string                 // 缓存名称
	version            *semver.Version        // 缓存版本
	indexActionCh      chan *core.IndexAction // 索引操作通道
	cachePostsCh       chan *postsEntry       // 推文缓存通道
	cache              tweetsCache            // 推文缓存接口
	lastCacheResetTime time.Time              // 上次缓存重置时间
	preventDuration    time.Duration          // 预防清除缓存的持续时间
}

// IndexPosts 根据用户、偏移和限制获取索引推文，优先从缓存获取，否则从数据库获取并更新缓存
func (s *cacheIndexSrv) IndexPosts(user *ms.User, offset int, limit int) (*ms.IndexTweetList, error) {
	key := s.keyFrom(user, offset, limit)
	posts, err := s.getPosts(key)
	if err == nil {
		logrus.Debugf("cacheIndexSrv.IndexPosts 通过键 %s 从缓存获取索引推文", key)
		return posts, nil
	}

	// 从数据库获取索引推文并更新缓存
	if posts, err = s.ips.IndexPosts(user, offset, limit); err != nil {
		return nil, err
	}
	logrus.Debugf("cacheIndexSrv.IndexPosts 通过键 %s 从数据库获取索引推文", key)
	s.cachePosts(key, posts)
	return posts, nil
}

// getPosts 从缓存获取推文数据
func (s *cacheIndexSrv) getPosts(key string) (*ms.IndexTweetList, error) {
	data, err := s.cache.getTweetsBytes(key)
	if err != nil {
		logrus.Debugf("cacheIndexSrv.getPosts 通过键 %s 从缓存获取推文时发生错误: %v", key, err)
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var resp ms.IndexTweetList
	if err := dec.Decode(&resp); err != nil {
		logrus.Debugf("cacheIndexSrv.getPosts 解码缓存中的推文时发生错误: %v", err)
		return nil, err
	}
	return &resp, nil
}

// cachePosts 将推文缓存到通道中
func (s *cacheIndexSrv) cachePosts(key string, tweets *ms.IndexTweetList) {
	entry := &postsEntry{key: key, tweets: tweets}
	select {
	case s.cachePostsCh <- entry:
		logrus.Debugf("cacheIndexSrv.cachePosts 通过通道缓存推文，键: %s", key)
	default:
		go func(ch chan<- *postsEntry, entry *postsEntry) {
			logrus.Debugf("cacheIndexSrv.cachePosts 通过goroutine缓存推文，键: %s", key)
			ch <- entry
		}(s.cachePostsCh, entry)
	}
}

// setPosts 将推文数据写入缓存
func (s *cacheIndexSrv) setPosts(entry *postsEntry) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(entry.tweets); err != nil {
		logrus.Debugf("cacheIndexSrv.setPosts 编码推文时发生错误: %v", err)
		return
	}
	if err := s.cache.setTweetsBytes(entry.key, buf.Bytes()); err != nil {
		logrus.Debugf("cacheIndexSrv.setPosts 设置缓存时发生错误: %v", err)
	}
	logrus.Debugf("cacheIndexSrv.setPosts 通过键 %s 设置缓存", entry.key)
}

// keyFrom 根据用户、偏移和限制生成缓存键
func (s *cacheIndexSrv) keyFrom(user *ms.User, offset int, limit int) string {
	var userId int64 = -1
	if user != nil {
		userId = user.ID
	}
	return fmt.Sprintf("%s:%d:%d:%d", _cacheIndexKey, userId, offset, limit)
}

// SendAction 发送索引操作到索引操作通道
func (s *cacheIndexSrv) SendAction(act core.IdxAct, post *ms.Post) {
	action := core.NewIndexAction(act, post)
	select {
	case s.indexActionCh <- action:
		logrus.Debugf("cacheIndexSrv.SendAction 通过通道发送索引操作: %s", act)
	default:
		go func(ch chan<- *core.IndexAction, act *core.IndexAction) {
			logrus.Debugf("cacheIndexSrv.SendAction 通过goroutine发送索引操作: %s", action.Act)
			ch <- act
		}(s.indexActionCh, action)
	}
}

// startIndexPosts 启动索引更新处理程序，处理推文缓存和索引操作
func (s *cacheIndexSrv) startIndexPosts() {
	for {
		select {
		case entry := <-s.cachePostsCh:
			s.setPosts(entry)
		case action := <-s.indexActionCh:
			s.handleIndexAction(action)
		}
	}
}

// handleIndexAction 处理索引操作，根据操作类型处理缓存清除
func (s *cacheIndexSrv) handleIndexAction(action *core.IndexAction) {
	act, post := action.Act, action.Post

	// 特殊处理创建/删除私密推文
	switch act {
	case core.IdxActCreatePost, core.IdxActDeletePost:
		if post.Visibility == core.PostVisitPrivate {
			s.deleteCacheByUserId(post.UserID, true)
			return
		}
	}

	// 如果在预防期内，清除所有缓存；否则只清除受影响的缓存（TODO：后续优化）
	if time.Since(s.lastCacheResetTime) > s.preventDuration {
		s.deleteCacheByUserId(post.UserID, false)
	} else {
		s.deleteCacheByUserId(post.UserID, true)
	}
}

// deleteCacheByUserId 根据用户ID删除缓存，可以选择是否只删除自身的缓存
func (s *cacheIndexSrv) deleteCacheByUserId(id int64, oneself bool) {
	var keys []string
	userId := strconv.FormatInt(id, 10)
	friendSet := ms.FriendSet{}
	if !oneself {
		friendSet = s.ams.MyFriendSet(id)
	}
	friendSet[userId] = types.Empty{}

	// 获取需要删除缓存的键，目前仅删除自己的缓存
	allKeys, err := s.cache.allKeys()
	if err != nil {
		logrus.Debugf("cacheIndexSrv.deleteCacheByUserId 用户ID: %s，发生错误：%s", userId, err)
	}
	for _, key := range allKeys {
		keyParts := strings.Split(key, ":")
		if len(keyParts) > 2 && keyParts[0] == _cacheIndexKey {
			if _, ok := friendSet[keyParts[1]]; ok {
				keys = append(keys, key)
			}
		}
	}

	// 执行删除缓存操作
	s.cache.delTweets(keys)
	s.lastCacheResetTime = time.Now()
	logrus.Debugf("cacheIndexSrv.deleteCacheByUserId 用户ID：%s 自身：%t 删除的键：%d", userId, oneself, len(keys))
}

// Name 返回缓存名称
func (s *cacheIndexSrv) Name() string {
	return s.name
}

// Version 返回缓存版本信息
func (s *cacheIndexSrv) Version() *semver.Version {
	return s.version
}

// newCacheIndexSrv 创建新的缓存索引服务实例
func newCacheIndexSrv(ips core.IndexPostsService, ams core.AuthorizationManageService, tc tweetsCache) *cacheIndexSrv {
	cacheIndex := &cacheIndexSrv{
		ips:             ips,
		ams:             ams,
		cache:           tc,
		name:            tc.Name(),
		version:         tc.Version(),
		preventDuration: 10 * time.Second,
	}

	// 根据配置文件中的最大更新QPS配置设置indexActionCh的容量，确保在 [10, 10000] 之间
	capacity := conf.CacheIndexSetting.MaxUpdateQPS
	if capacity < 10 {
		capacity = 10
	} else if capacity > 10000 {
		capacity = 10000
	}
	cacheIndex.indexActionCh = make(chan *core.IndexAction, capacity)
	cacheIndex.cachePostsCh = make(chan *postsEntry, capacity)
	// 启动索引更新处理程序
	go cacheIndex.startIndexPosts()

	return cacheIndex
}

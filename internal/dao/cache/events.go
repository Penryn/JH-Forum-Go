package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/alimy/tryst/event"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/events"
	"github.com/sirupsen/logrus"
)

// BaseCacheEvent 定义基础的缓存事件结构
type BaseCacheEvent struct {
	event.UnimplementedEvent               // 嵌入事件接口
	ac                       core.AppCache // 应用缓存实例
}

// expireIndexTweetsEvent 定义过期索引推文事件结构
type expireIndexTweetsEvent struct {
	event.UnimplementedEvent               // 嵌入事件接口
	ac                       core.AppCache // 应用缓存实例
	keysPattern              []string      // 缓存键模式列表
}

// expireHotsTweetsEvent 定义过期热门推文事件结构
type expireHotsTweetsEvent struct {
	event.UnimplementedEvent               // 嵌入事件接口
	ac                       core.AppCache // 应用缓存实例
	keyPattern               string        // 缓存键模式
}

// expireFollowTweetsEvent 定义过期关注推文事件结构
type expireFollowTweetsEvent struct {
	event.UnimplementedEvent               // 嵌入事件接口
	tweet                    *ms.Post      // 推文对象
	ac                       core.AppCache // 应用缓存实例
	keyPattern               string        // 缓存键模式
}

// cacheObjectEvent 定义缓存对象事件结构
type cacheObjectEvent struct {
	event.UnimplementedEvent               // 嵌入事件接口
	ac                       core.AppCache // 应用缓存实例
	key                      string        // 缓存键
	data                     any           // 缓存数据
	expire                   int64         // 过期时间
}

// cacheUserInfoEvent 定义缓存用户信息事件结构
type cacheUserInfoEvent struct {
	event.UnimplementedEvent               // 嵌入事件接口
	ac                       core.AppCache // 应用缓存实例
	key                      string        // 缓存键
	data                     *ms.User      // 用户信息数据
	expire                   int64         // 过期时间
}

// cacheMyFriendIdsEvent 定义缓存我的好友ID事件结构
type cacheMyFriendIdsEvent struct {
	event.UnimplementedEvent                          // 嵌入事件接口
	ac                       core.AppCache            // 应用缓存实例
	urs                      core.UserRelationService // 用户关系服务接口
	userIds                  []int64                  // 用户ID列表
	expire                   int64                    // 过期时间
}

// cacheMyFollowIdsEvent 定义缓存我的关注ID事件结构
type cacheMyFollowIdsEvent struct {
	event.UnimplementedEvent                          // 嵌入事件接口
	ac                       core.AppCache            // 应用缓存实例
	urs                      core.UserRelationService // 用户关系服务接口
	userId                   int64                    // 用户ID
	key                      string                   // 缓存键
	expire                   int64                    // 过期时间
}

// NewBaseCacheEvent 创建一个新的基础缓存事件实例
func NewBaseCacheEvent(ac core.AppCache) *BaseCacheEvent {
	return &BaseCacheEvent{
		ac: ac,
	}
}

// OnExpireIndexTweetEvent 触发过期索引推文事件，使相关缓存过期
func OnExpireIndexTweetEvent(userId int64) {
	events.OnEvent(&expireIndexTweetsEvent{
		ac: _appCache,
		keysPattern: []string{
			conf.PrefixIdxTweetsNewest + "*",
			conf.PrefixIdxTweetsHots + "*",
			conf.PrefixIdxTweetsFollowing + "*",
			fmt.Sprintf("%s%d:*", conf.PrefixUserTweets, userId),
		},
	})
}

// OnExpireHotsTweetEvent 触发过期热门推文事件，使相关缓存过期
func OnExpireHotsTweetEvent() {
	events.OnEvent(&expireHotsTweetsEvent{
		ac:         _appCache,
		keyPattern: conf.PrefixHotsTweets + "*",
	})
}

// onExpireFollowTweetEvent 触发过期关注推文事件，使相关缓存过期
func onExpireFollowTweetEvent(tweet *ms.Post) {
	events.OnEvent(&expireFollowTweetsEvent{
		tweet:      tweet,
		ac:         _appCache,
		keyPattern: conf.PrefixFollowingTweets + "*",
	})
}

// onCacheUserInfoEvent 触发缓存用户信息事件，更新缓存
func onCacheUserInfoEvent(key string, data *ms.User) {
	events.OnEvent(&cacheUserInfoEvent{
		key:    key,
		data:   data,
		ac:     _appCache,
		expire: conf.CacheSetting.UserInfoExpire,
	})
}

// onCacheObjectEvent 触发缓存对象事件，更新缓存
func onCacheObjectEvent(key string, data any, expire int64) {
	events.OnEvent(&cacheObjectEvent{
		key:    key,
		data:   data,
		ac:     _appCache,
		expire: expire,
	})
}

// OnCacheMyFriendIdsEvent 触发缓存我的好友ID事件，更新缓存
func OnCacheMyFriendIdsEvent(urs core.UserRelationService, userIds ...int64) {
	if len(userIds) == 0 {
		return
	}
	events.OnEvent(&cacheMyFriendIdsEvent{
		userIds: userIds,
		urs:     urs,
		ac:      _appCache,
		expire:  conf.CacheSetting.UserRelationExpire,
	})
}

// OnCacheMyFollowIdsEvent 触发缓存我的关注ID事件，更新缓存
func OnCacheMyFollowIdsEvent(urs core.UserRelationService, userId int64, key ...string) {
	cacheKey := ""
	if len(key) > 0 {
		cacheKey = key[0]
	} else {
		cacheKey = conf.KeyMyFollowIds.Get(userId)
	}
	events.OnEvent(&cacheMyFollowIdsEvent{
		userId: userId,
		urs:    urs,
		key:    cacheKey,
		ac:     _appCache,
		expire: conf.CacheSetting.UserRelationExpire,
	})
}

// ExpireUserInfo 根据用户ID和用户名使用户信息缓存过期
func (e *BaseCacheEvent) ExpireUserInfo(id int64, name string) error {
	keys := make([]string, 0, 2)
	if id >= 0 {
		keys = append(keys, conf.KeyUserInfoById.Get(id))
	}
	if len(name) > 0 {
		keys = append(keys, conf.KeyUserInfoByName.Get(name))
	}
	return e.ac.Delete(keys...)
}

// ExpireUserProfile 根据用户名使用户资料缓存过期
func (e *BaseCacheEvent) ExpireUserProfile(name string) error {
	if len(name) > 0 {
		return e.ac.Delete(conf.KeyUserProfileByName.Get(name))
	}
	return nil
}

// ExpireUserData 根据用户ID和用户名使用户相关数据缓存过期
func (e *BaseCacheEvent) ExpireUserData(id int64, name string) error {
	keys := make([]string, 0, 3)
	if id >= 0 {
		keys = append(keys, conf.KeyUserInfoById.Get(id))
	}
	if len(name) > 0 {
		keys = append(keys, conf.KeyUserInfoByName.Get(name), conf.KeyUserProfileByName.Get(name))
	}
	return e.ac.Delete(keys...)
}

// Name 返回过期索引推文事件的名称
func (e *expireIndexTweetsEvent) Name() string {
	return "expireIndexTweetsEvent"
}

// Action 执行过期索引推文事件的操作
func (e *expireIndexTweetsEvent) Action() (err error) {
	for _, pattern := range e.keysPattern {
		e.ac.DelAny(pattern)
	}
	return
}

// Name 返回过期热门推文事件的名称
func (e *expireHotsTweetsEvent) Name() string {
	return "expireHotsTweetsEvent"
}

// Action 执行过期热门推文事件的操作
func (e *expireHotsTweetsEvent) Action() (err error) {
	e.ac.DelAny(e.keyPattern)
	return
}

// Name 返回过期关注推文事件的名称
func (e *expireFollowTweetsEvent) Name() string {
	return "expireFollowTweetsEvent"
}

// Action 执行过期关注推文事件的操作
func (e *expireFollowTweetsEvent) Action() (err error) {
	e.ac.DelAny(e.keyPattern)
	return
}

// Name 返回缓存用户信息事件的名称
func (e *cacheUserInfoEvent) Name() string {
	return "cacheUserInfoEvent"
}

// Action 执行缓存用户信息事件的操作
func (e *cacheUserInfoEvent) Action() (err error) {
	buffer := &bytes.Buffer{}
	ge := gob.NewEncoder(buffer)
	if err = ge.Encode(e.data); err == nil {
		e.ac.Set(e.key, buffer.Bytes(), e.expire)
	}
	return
}

// Name 返回缓存对象事件的名称
func (e *cacheObjectEvent) Name() string {
	return "cacheObjectEvent"
}

// Action 执行缓存对象事件的操作
func (e *cacheObjectEvent) Action() (err error) {
	buffer := &bytes.Buffer{}
	ge := gob.NewEncoder(buffer)
	if err = ge.Encode(e.data); err == nil {
		e.ac.Set(e.key, buffer.Bytes(), e.expire)
	}
	return
}

// Name 返回缓存我的好友ID事件的名称
func (e *cacheMyFriendIdsEvent) Name() string {
	return "cacheMyFriendIdsEvent"
}

// Action 执行缓存我的好友ID事件的操作
func (e *cacheMyFriendIdsEvent) Action() error {
	logrus.Debug("cacheMyFriendIdsEvent action running")
	for _, userId := range e.userIds {
		myFriendIds, err := e.urs.MyFriendIds(userId)
		if err != nil {
			return err
		}
		bitmap := roaring64.New()
		for _, friendId := range myFriendIds {
			bitmap.Add(uint64(friendId))
		}
		data, err := bitmap.MarshalBinary()
		if err != nil {
			return err
		}
		e.ac.Set(conf.KeyMyFriendIds.Get(userId), data, e.expire)
	}
	return nil
}

// Name 返回缓存我的关注ID事件的名称
func (e *cacheMyFollowIdsEvent) Name() string {
	return "cacheMyFollowIdsEvent"
}

// Action 执行缓存我的关注ID事件的操作
func (e *cacheMyFollowIdsEvent) Action() (err error) {
	logrus.Debug("cacheMyFollowIdsEvent action running")
	myFollowIds, err := e.urs.MyFollowIds(e.userId)
	if err != nil {
		return err
	}
	bitmap := roaring64.New()
	for _, followId := range myFollowIds {
		bitmap.Add(uint64(followId))
	}
	data, err := bitmap.MarshalBinary()
	if err != nil {
		return err
	}
	e.ac.Set(e.key, data, e.expire)
	return nil
}

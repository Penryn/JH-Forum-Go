// 该代码文件定义了在缓存中使用的一些键前缀和池化后的键。

package conf

import (
	"fmt"

	"github.com/alimy/tryst/cache"
	"JH-Forum/pkg/types"
)

const (
	_defaultKeyPoolSize = 128
)

// 以下包含一些在缓存中会用到的键前缀
const (
	InfixCommentDefault      = "default"
	InfixCommentHots         = "hots"
	InfixCommentNewest       = "newest"
	PrefixNewestTweets       = "paopao:newesttweets:"
	PrefixHotsTweets         = "paopao:hotstweets:"
	PrefixFollowingTweets    = "paopao:followingtweets:"
	PrefixUserTweets         = "paopao:usertweets:"
	PrefixUnreadmsg          = "paopao:unreadmsg:"
	PrefixOnlineUser         = "paopao:onlineuser:"
	PrefixIdxTweetsNewest    = "paopao:index:tweets:newest:"
	PrefixIdxTweetsHots      = "paopao:index:tweets:hots:"
	PrefixIdxTweetsFollowing = "paopao:index:tweets:following:"
	PrefixIdxTrends          = "paopao:index:trends:"
	PrefixMessages           = "paopao:messages:"
	PrefixUserInfo           = "paopao:user:info:"
	PrefixUserProfile        = "paopao:user:profile:"
	PrefixUserInfoById       = "paopao:user:info:id:"
	PrefixUserInfoByName     = "paopao:user:info:name:"
	prefixUserProfileByName  = "paopao:user:profile:name:"
	PrefixMyFriendIds        = "paopao:myfriendids:"
	PrefixMyFollowIds        = "paopao:myfollowids:"
	PrefixTweetComment       = "paopao:comment:"
	KeySiteStatus            = "paopao:sitestatus"
	KeyHistoryMaxOnline      = "history.max.online"
)

// 以下包含一些在缓存中会用到的池化后的键
var (
	KeyNewestTweets      cache.KeyPool[int]    // 用于存储新推文的池化键
	KeyHotsTweets        cache.KeyPool[int]    // 用于存储热门推文的池化键
	KeyFollowingTweets   cache.KeyPool[string] // 用于存储关注用户推文的池化键
	KeyUnreadMsg         cache.KeyPool[int64]  // 用于存储未读消息的池化键
	KeyOnlineUser        cache.KeyPool[int64]  // 用于存储在线用户的池化键
	KeyUserInfoById      cache.KeyPool[int64]  // 用于存储用户信息按ID索引的池化键
	KeyUserInfoByName    cache.KeyPool[string] // 用于存储用户信息按用户名索引的池化键
	KeyUserProfileByName cache.KeyPool[string] // 用于存储用户配置信息按用户名索引的池化键
	KeyMyFriendIds       cache.KeyPool[int64]  // 用于存储我的好友ID列表的池化键
	KeyMyFollowIds       cache.KeyPool[int64]  // 用于存储我关注的用户ID列表的池化键
)

// initCacheKeyPool 初始化缓存键池
func initCacheKeyPool() {
	poolSize := _defaultKeyPoolSize
	if poolSize < CacheSetting.KeyPoolSize {
		poolSize = CacheSetting.KeyPoolSize
	}
	// 初始化各个池化键
	KeyNewestTweets = intKeyPool[int](poolSize, PrefixNewestTweets)
	KeyHotsTweets = intKeyPool[int](poolSize, PrefixHotsTweets)
	KeyFollowingTweets = strKeyPool(poolSize, PrefixFollowingTweets)
	KeyUnreadMsg = intKeyPool[int64](poolSize, PrefixUnreadmsg)
	KeyOnlineUser = intKeyPool[int64](poolSize, PrefixOnlineUser)
	KeyUserInfoById = intKeyPool[int64](poolSize, PrefixUserInfoById)
	KeyUserInfoByName = strKeyPool(poolSize, PrefixUserInfoByName)
	KeyUserProfileByName = strKeyPool(poolSize, prefixUserProfileByName)
	KeyMyFriendIds = intKeyPool[int64](poolSize, PrefixMyFriendIds)
	KeyMyFollowIds = intKeyPool[int64](poolSize, PrefixMyFollowIds)
}

// strKeyPool 创建字符串类型的缓存键池
func strKeyPool(size int, prefix string) cache.KeyPool[string] {
	return cache.MustKeyPool(size, func(key string) string {
		return fmt.Sprintf("%s%s", prefix, key)
	})
}

// intKeyPool 创建整数类型的缓存键池
func intKeyPool[T types.Integer](size int, prefix string) cache.KeyPool[T] {
	return cache.MustKeyPool[T](size, intKey[T](prefix))
}

// intKey 创建整数类型的键生成函数
func intKey[T types.Integer](prefix string) func(T) string {
	return func(key T) string {
		return fmt.Sprintf("%s%d", prefix, key)
	}
}

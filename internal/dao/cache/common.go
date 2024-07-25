package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/RoaringBitmap/roaring/roaring64"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
)

// cacheDataService 是一个实现了缓存功能的数据服务
type cacheDataService struct {
	core.DataService               // 嵌入核心数据服务接口
	ac               core.AppCache // 应用缓存实例
}

// NewCacheDataService 返回一个包装了缓存功能的数据服务实例
func NewCacheDataService(ds core.DataService) core.DataService {
	lazyInitial() // 懒加载应用缓存
	return &cacheDataService{
		DataService: ds,
		ac:          _appCache,
	}
}

// GetUserByID 通过用户ID获取用户信息，优先从缓存获取，缓存命中后更新缓存
func (s *cacheDataService) GetUserByID(id int64) (res *ms.User, err error) {
	key := conf.KeyUserInfoById.Get(id)           // 获取缓存键
	if data, xerr := s.ac.Get(key); xerr == nil { // 尝试从缓存获取数据
		buf := bytes.NewBuffer(data)
		res = &ms.User{}
		err = gob.NewDecoder(buf).Decode(res)
		return
	}
	// 缓存未命中，从数据库获取
	if res, err = s.DataService.GetUserByID(id); err == nil {
		onCacheUserInfoEvent(key, res) // 更新缓存
	}
	return
}

// GetUserByUsername 通过用户名获取用户信息，优先从缓存获取，缓存命中后更新缓存
func (s *cacheDataService) GetUserByUsername(username string) (res *ms.User, err error) {
	key := conf.KeyUserInfoByName.Get(username)   // 获取缓存键
	if data, xerr := s.ac.Get(key); xerr == nil { // 尝试从缓存获取数据
		buf := bytes.NewBuffer(data)
		res = &ms.User{}
		err = gob.NewDecoder(buf).Decode(res)
		return
	}
	// 缓存未命中，从数据库获取
	if res, err = s.DataService.GetUserByUsername(username); err == nil {
		onCacheUserInfoEvent(key, res) // 更新缓存
	}
	return
}

// UserProfileByName 通过用户名获取用户资料，优先从缓存获取，缓存命中后更新缓存
func (s *cacheDataService) UserProfileByName(username string) (res *cs.UserProfile, err error) {
	key := conf.KeyUserProfileByName.Get(username) // 获取缓存键
	if data, xerr := s.ac.Get(key); xerr == nil {  // 尝试从缓存获取数据
		buf := bytes.NewBuffer(data)
		res = &cs.UserProfile{}
		err = gob.NewDecoder(buf).Decode(res)
		return
	}
	// 缓存未命中，从数据库获取
	if res, err = s.DataService.UserProfileByName(username); err == nil {
		onCacheObjectEvent(key, res, conf.CacheSetting.UserProfileExpire) // 更新缓存
	}
	return
}

// IsMyFriend 检查是否为好友关系，优先从缓存获取，缓存命中后更新缓存
func (s *cacheDataService) IsMyFriend(userId int64, friendIds ...int64) (res map[int64]bool, err error) {
	size := len(friendIds)
	res = make(map[int64]bool, size)
	if size == 0 {
		return
	}
	key := conf.KeyMyFriendIds.Get(userId)        // 获取缓存键
	if data, xerr := s.ac.Get(key); xerr == nil { // 尝试从缓存获取数据
		bitmap := roaring64.New()
		if err = bitmap.UnmarshalBinary(data); err == nil {
			for _, friendId := range friendIds {
				res[friendId] = bitmap.Contains(uint64(friendId))
			}
			return
		}
	}
	// 缓存未命中，直接从数据库获取并触发缓存更新事件
	OnCacheMyFriendIdsEvent(s.DataService, userId)
	return s.DataService.IsMyFriend(userId, friendIds...)
}

// IsMyFollow 检查是否为关注关系，优先从缓存获取，缓存命中后更新缓存
func (s *cacheDataService) IsMyFollow(userId int64, followIds ...int64) (res map[int64]bool, err error) {
	size := len(followIds)
	res = make(map[int64]bool, size)
	if size == 0 {
		return
	}
	key := conf.KeyMyFollowIds.Get(userId)        // 获取缓存键
	if data, xerr := s.ac.Get(key); xerr == nil { // 尝试从缓存获取数据
		bitmap := roaring64.New()
		if err = bitmap.UnmarshalBinary(data); err == nil {
			for _, followId := range followIds {
				res[followId] = bitmap.Contains(uint64(followId))
			}
			return
		}
	}
	// 缓存未命中，直接从数据库获取并触发缓存更新事件
	OnCacheMyFollowIdsEvent(s.DataService, userId, key)
	return s.DataService.IsMyFollow(userId, followIds...)
}

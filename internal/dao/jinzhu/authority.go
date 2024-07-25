package jinzhu

import (
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/dao/jinzhu/dbr"
	"JH-Forum/pkg/types"
	"gorm.io/gorm"
)

// authorizationManageSrv 实现了 AuthorizationManageService 接口，用于处理授权管理相关逻辑
type authorizationManageSrv struct {
	db *gorm.DB
}

// newAuthorizationManageService 创建一个新的 authorizationManageSrv 实例
func newAuthorizationManageService(db *gorm.DB) core.AuthorizationManageService {
	return &authorizationManageSrv{
		db: db,
	}
}

// IsAllow 判断用户是否有权限执行指定动作
func (s *authorizationManageSrv) IsAllow(user *ms.User, action *ms.Action) bool {
	// 用户已激活，如果已绑定手机号
	isActivation := (len(user.Phone) != 0)
	isFriend := s.isFriend(user.ID, action.UserId)
	// TODO: 目前仅使用默认的授权检查规则
	return action.Act.IsAllow(user, action.UserId, isFriend, isActivation)
}

// MyFriendSet 获取指定用户的好友集合
func (s *authorizationManageSrv) MyFriendSet(userId int64) ms.FriendSet {
	ids, err := (&dbr.Contact{UserId: userId}).MyFriendIds(s.db)
	if err != nil {
		return ms.FriendSet{}
	}

	resp := make(ms.FriendSet, len(ids))
	for _, id := range ids {
		resp[id] = types.Empty{}
	}
	return resp
}

// BeFriendFilter 获取指定用户作为被加好友的过滤器
func (s *authorizationManageSrv) BeFriendFilter(userId int64) ms.FriendFilter {
	ids, err := (&dbr.Contact{FriendId: userId}).BeFriendIds(s.db)
	if err != nil {
		return ms.FriendFilter{}
	}

	resp := make(ms.FriendFilter, len(ids))
	for _, id := range ids {
		resp[id] = types.Empty{}
	}
	return resp
}

// BeFriendIds 获取指定用户的好友列表
func (s *authorizationManageSrv) BeFriendIds(userId int64) ([]int64, error) {
	return (&dbr.Contact{FriendId: userId}).BeFriendIds(s.db)
}

// isFriend 检查两个用户是否为好友关系
func (s *authorizationManageSrv) isFriend(userId int64, friendId int64) bool {
	contact, err := (&dbr.Contact{UserId: friendId, FriendId: userId}).GetByUserFriend(s.db)
	if err == nil || contact.Status == dbr.ContactStatusAgree {
		return true
	}
	return false
}

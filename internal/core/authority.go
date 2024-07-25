// 该接口定义了授权管理服务的功能，包括权限检查、好友过滤和好友集管理。

package core

import (
	"JH-Forum/internal/core/ms"
)

// AuthorizationManageService 授权管理服务
type AuthorizationManageService interface {
	// 检查用户是否有权限执行某个操作
	IsAllow(user *ms.User, action *ms.Action) bool

	// 获取用户的好友过滤器
	BeFriendFilter(userId int64) ms.FriendFilter

	// 获取用户的好友ID列表
	BeFriendIds(userId int64) ([]int64, error)

	// 获取用户的好友集
	MyFriendSet(userId int64) ms.FriendSet
}

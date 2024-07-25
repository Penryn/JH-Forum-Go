// 该代码定义了多个接口，用于实现用户管理、联系人管理、关注管理和用户关系服务。

package core

import (
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
)

// UserManageService 用户管理服务接口
type UserManageService interface {
	GetUserByID(id int64) (*ms.User, error)                     // 通过ID获取用户
	GetUserByUsername(username string) (*ms.User, error)        // 通过用户名获取用户
	GetUserByPhone(phone string) (*ms.User, error)              // 通过手机号获取用户
	GetUsersByIDs(ids []int64) ([]*ms.User, error)              // 通过ID列表获取用户
	GetUsersByKeyword(keyword string) ([]*ms.User, error)       // 通过关键词获取用户
	UserProfileByName(username string) (*cs.UserProfile, error) // 通过用户名获取用户资料
	CreateUser(user *ms.User) (*ms.User, error)                 // 创建用户
	UpdateUser(user *ms.User) error                             // 更新用户信息
	GetRegisterUserCount() (int64, error)                       // 获取注册用户数量
}

// ContactManageService 联系人管理服务接口
type ContactManageService interface {
	RequestingFriend(userId int64, friendId int64, greetings string) error    // 申请添加好友
	AddFriend(userId int64, friendId int64) error                             // 添加好友
	RejectFriend(userId int64, friendId int64) error                          // 拒绝好友请求
	DeleteFriend(userId int64, friendId int64) error                          // 删除好友
	GetContacts(userId int64, offset int, limit int) (*ms.ContactList, error) // 获取联系人列表
	IsFriend(userID int64, friendID int64) bool                               // 判断是否为好友
}

// FollowingManageService 关注管理服务接口
type FollowingManageService interface {
	FollowUser(userId int64, followId int64) error                           // 关注用户
	UnfollowUser(userId int64, followId int64) error                         // 取消关注用户
	ListFollows(userId int64, limit, offset int) (*ms.ContactList, error)    // 列出关注者
	ListFollowings(userId int64, limit, offset int) (*ms.ContactList, error) // 列出关注的用户
	GetFollowCount(userId int64) (int64, int64, error)                       // 获取关注数量
	IsFollow(userId int64, followId int64) bool                              // 判断是否关注
}

// UserRelationService 用户关系服务接口
type UserRelationService interface {
	MyFriendIds(userId int64) ([]int64, error)                           // 获取我的好友ID列表
	MyFollowIds(userId int64) ([]int64, error)                           // 获取我的关注ID列表
	IsMyFriend(userId int64, friendIds ...int64) (map[int64]bool, error) // 判断是否是我的好友
	IsMyFollow(userId int64, followIds ...int64) (map[int64]bool, error) // 判断是否是我的关注
}

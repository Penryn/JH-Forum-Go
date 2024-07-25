// Package jinzhu 实现了用户管理和用户关系服务。
package jinzhu

import (
	"fmt"
	"strings"

	"JH-Forum/internal/core"
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/dao/jinzhu/dbr"
	"gorm.io/gorm"
)

// userManageSrv 实现了 core.UserManageService 接口，提供用户管理的各种功能。
type userManageSrv struct {
	db  *gorm.DB
	ums core.UserMetricServantA

	_userProfileJoins    string
	_userProfileWhere    string
	_userProfileColoumns []string
}

// userRelationSrv 实现了用户关系服务，包括获取好友和关注列表等功能。
type userRelationSrv struct {
	db *gorm.DB
}

// newUserManageService 创建并返回一个新的 userManageSrv 实例。
func newUserManageService(db *gorm.DB, ums core.UserMetricServantA) core.UserManageService {
	return &userManageSrv{
		db:                db,
		ums:               ums,
		_userProfileJoins: fmt.Sprintf("LEFT JOIN %s m ON %s.id=m.user_id", _userMetric_, _user_),
		_userProfileWhere: fmt.Sprintf("%s.username=? AND %s.is_del=0", _user_, _user_),
		_userProfileColoumns: []string{
			fmt.Sprintf("%s.id", _user_),
			fmt.Sprintf("%s.username", _user_),
			fmt.Sprintf("%s.nickname", _user_),
			fmt.Sprintf("%s.phone", _user_),
			fmt.Sprintf("%s.status", _user_),
			fmt.Sprintf("%s.avatar", _user_),
			fmt.Sprintf("%s.balance", _user_),
			fmt.Sprintf("%s.is_admin", _user_),
			fmt.Sprintf("%s.created_on", _user_),
			"m.tweets_count",
		},
	}
}

// newUserRelationService 创建并返回一个新的 userRelationSrv 实例。
func newUserRelationService(db *gorm.DB) core.UserRelationService {
	return &userRelationSrv{
		db: db,
	}
}

// GetUserByID 根据用户 ID 获取用户信息。
func (s *userManageSrv) GetUserByID(id int64) (*ms.User, error) {
	user := &dbr.User{
		Model: &dbr.Model{
			ID: id,
		},
	}
	return user.Get(s.db)
}

// GetUserByUsername 根据用户名获取用户信息。
func (s *userManageSrv) GetUserByUsername(username string) (*ms.User, error) {
	user := &dbr.User{
		Username: username,
	}
	return user.Get(s.db)
}

// UserProfileByName 根据用户名获取用户的详细信息。
func (s *userManageSrv) UserProfileByName(username string) (res *cs.UserProfile, err error) {
	err = s.db.Table(_user_).Joins(s._userProfileJoins).
		Where(s._userProfileWhere, username).
		Select(s._userProfileColoumns).
		First(&res).Error
	return
}

// GetUserByPhone 根据手机号获取用户信息。
func (s *userManageSrv) GetUserByPhone(phone string) (*ms.User, error) {
	user := &dbr.User{
		Phone: phone,
	}
	return user.Get(s.db)
}

// GetUsersByIDs 根据用户 ID 列表获取用户信息列表。
func (s *userManageSrv) GetUsersByIDs(ids []int64) ([]*ms.User, error) {
	user := &dbr.User{}
	return user.List(s.db, &dbr.ConditionsT{
		"id IN ?": ids,
	}, 0, 0)
}

// GetUsersByKeyword 根据关键词搜索用户。
func (s *userManageSrv) GetUsersByKeyword(keyword string) ([]*ms.User, error) {
	user := &dbr.User{}
	keyword = strings.Trim(keyword, " ") + "%"
	if keyword == "%" {
		return user.List(s.db, &dbr.ConditionsT{
			"ORDER": "id ASC",
		}, 0, 6)
	} else {
		return user.List(s.db, &dbr.ConditionsT{
			"username LIKE ?": keyword,
		}, 0, 6)
	}
}

// CreateUser 创建新用户。
func (s *userManageSrv) CreateUser(user *dbr.User) (res *ms.User, err error) {
	if res, err = user.Create(s.db); err == nil {
		s.ums.AddUserMetric(res.ID)
	}
	return
}

// UpdateUser 更新用户信息。
func (s *userManageSrv) UpdateUser(user *ms.User) error {
	return user.Update(s.db)
}

// GetRegisterUserCount 获取注册用户数量。
func (s *userManageSrv) GetRegisterUserCount() (res int64, err error) {
	err = s.db.Model(&dbr.User{}).Count(&res).Error
	return
}

// MyFriendIds 获取用户的好友 ID 列表。
func (s *userRelationSrv) MyFriendIds(userId int64) (res []int64, err error) {
	err = s.db.Table(_contact_).Where("user_id=? AND status=2 AND is_del=0", userId).Select("friend_id").Find(&res).Error
	return
}

// MyFollowIds 获取用户的关注 ID 列表。
func (s *userRelationSrv) MyFollowIds(userId int64) (res []int64, err error) {
	err = s.db.Table(_following_).Where("user_id=? AND is_del=0", userId).Select("follow_id").Find(&res).Error
	return
}

// IsMyFriend 检查用户是否是指定用户的好友。
func (s *userRelationSrv) IsMyFriend(userId int64, friendIds ...int64) (map[int64]bool, error) {
	size := len(friendIds)
	res := make(map[int64]bool, size)
	if size == 0 {
		return res, nil
	}
	myFriendIds, err := s.MyFriendIds(userId)
	if err != nil {
		return nil, err
	}
	for _, friendId := range friendIds {
		res[friendId] = false
		for _, myFriendId := range myFriendIds {
			if friendId == myFriendId {
				res[friendId] = true
				break
			}
		}
	}
	return res, nil
}

// IsMyFollow 检查用户是否关注了指定用户。
func (s *userRelationSrv) IsMyFollow(userId int64, followIds ...int64) (map[int64]bool, error) {
	size := len(followIds)
	res := make(map[int64]bool, size)
	if size == 0 {
		return res, nil
	}
	myFollowIds, err := s.MyFollowIds(userId)
	if err != nil {
		return nil, err
	}
	for _, followId := range followIds {
		res[followId] = false
		for _, myFollowId := range myFollowIds {
			if followId == myFollowId {
				res[followId] = true
				break
			}
		}
	}
	return res, nil
}

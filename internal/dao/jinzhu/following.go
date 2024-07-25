package jinzhu

import (
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/ms"
	"JH-Forum/internal/dao/jinzhu/dbr"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	_ core.FollowingManageService = (*followingManageSrv)(nil)
)

type followingManageSrv struct {
	db *gorm.DB
	f  *dbr.Following
	u  *dbr.User
}

// 创建新的关注管理服务
func newFollowingManageService(db *gorm.DB) core.FollowingManageService {
	return &followingManageSrv{
		db: db,
		f:  &dbr.Following{},
		u:  &dbr.User{},
	}
}

// 关注用户
func (s *followingManageSrv) FollowUser(userId int64, followId int64) error {
	if _, err := s.f.GetFollowing(s.db, userId, followId); err != nil {
		following := &dbr.Following{
			UserId:   userId,
			FollowId: followId,
		}
		if _, err = following.Create(s.db); err != nil {
			logrus.Errorf("contactManageSrv.fetchOrNewContact create new contact err:%s", err)
			return err
		}
	}
	return nil
}

// 取消关注用户
func (s *followingManageSrv) UnfollowUser(userId int64, followId int64) error {
	return s.f.DelFollowing(s.db, userId, followId)
}

// 获取用户关注列表
func (s *followingManageSrv) ListFollows(userId int64, limit, offset int) (*ms.ContactList, error) {
	follows, total, err := s.f.ListFollows(s.db, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	res := &ms.ContactList{
		Total: total,
	}
	for _, f := range follows {
		res.Contacts = append(res.Contacts, ms.ContactItem{
			UserId:    f.User.ID,
			Username:  f.User.Username,
			Nickname:  f.User.Nickname,
			Avatar:    f.User.Avatar,
			CreatedOn: f.User.CreatedOn,
		})
	}
	return res, nil
}

// 获取用户粉丝列表
func (s *followingManageSrv) ListFollowings(userId int64, limit, offset int) (*ms.ContactList, error) {
	followingIds, total, err := s.f.ListFollowingIds(s.db, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	followings, err := s.u.ListUserInfoById(s.db, followingIds)
	if err != nil {
		return nil, err
	}
	res := &ms.ContactList{
		Total: total,
	}
	for _, user := range followings {
		res.Contacts = append(res.Contacts, ms.ContactItem{
			UserId:    user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			CreatedOn: user.CreatedOn,
		})
	}
	return res, nil
}

// 获取用户关注数和粉丝数
func (s *followingManageSrv) GetFollowCount(userId int64) (int64, int64, error) {
	return s.f.FollowCount(s.db, userId)
}

// 判断是否关注了用户
func (s *followingManageSrv) IsFollow(userId int64, followId int64) bool {
	if _, err := s.f.GetFollowing(s.db, userId, followId); err == nil {
		return true
	}
	return false
}

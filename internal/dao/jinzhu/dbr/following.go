package dbr

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Following 表示用户关注关系的对象，包括用户ID、关注者ID。
type Following struct {
	*Model         // 模型基类
	User     *User `json:"-" gorm:"foreignKey:ID;references:FollowId"` // 用户对象
	UserId   int64 `json:"user_id"`                                    // 用户ID
	FollowId int64 `json:"friend_id"`                                  // 关注者ID
}

// GetFollowing 方法从数据库中获取指定用户关注的记录。
func (f *Following) GetFollowing(db *gorm.DB, userId, followId int64) (*Following, error) {
	var following Following
	err := db.Omit("User").Unscoped().Where("user_id = ? AND follow_id = ?", userId, followId).First(&following).Error
	if err != nil {
		logrus.Debugf("Following.GetFollowing get following error:%s", err)
		return nil, err
	}
	return &following, nil
}

// DelFollowing 方法从数据库中删除指定用户关注的记录。
func (f *Following) DelFollowing(db *gorm.DB, userId, followId int64) error {
	return db.Omit("User").Unscoped().Where("user_id = ? AND follow_id = ?", userId, followId).Delete(f).Error
}

// ListFollows 方法从数据库中获取指定用户关注的所有记录。
func (f *Following) ListFollows(db *gorm.DB, userId int64, limit int, offset int) (res []*Following, total int64, err error) {
	db = db.Model(f).Where("user_id=?", userId)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	db.Joins("User").Order(clause.OrderByColumn{Column: clause.Column{Table: "User", Name: "nickname"}, Desc: false})
	if err = db.Find(&res).Error; err != nil {
		return
	}
	return
}

// ListFollowingIds 方法从数据库中获取指定用户被关注的所有用户ID。
func (f *Following) ListFollowingIds(db *gorm.DB, userId int64, limit, offset int) (ids []int64, total int64, err error) {
	db = db.Model(f).Where("follow_id=?", userId)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Omit("User").Select("user_id").Find(&ids).Error; err != nil {
		return
	}
	return
}

// FollowCount 方法统计指定用户的关注和粉丝数量。
func (f *Following) FollowCount(db *gorm.DB, userId int64) (follows int64, followings int64, err error) {
	if err = db.Model(f).Where("user_id=?", userId).Count(&follows).Error; err != nil {
		return
	}
	if err = db.Model(f).Where("follow_id=?", userId).Count(&followings).Error; err != nil {
		return
	}
	return
}

// IsFollow 方法检查指定用户是否关注了另一个用户。
func (s *Following) IsFollow(db *gorm.DB, userId int64, followId int64) bool {
	if _, err := s.GetFollowing(db, userId, followId); err == nil {
		return true
	}
	return false
}

// Create 方法在数据库中创建新的关注关系记录。
func (f *Following) Create(db *gorm.DB) (*Following, error) {
	err := db.Omit("User").Create(f).Error
	return f, err
}

// UpdateInUnscoped 方法在不考虑软删除的情况下更新数据库中的关注关系记录。
func (c *Following) UpdateInUnscoped(db *gorm.DB) error {
	return db.Unscoped().Omit("User").Save(c).Error
}

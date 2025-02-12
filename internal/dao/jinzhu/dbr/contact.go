package dbr

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ContactStatusRequesting int8 = iota + 1
	ContactStatusAgree
	ContactStatusReject
	ContactStatusDeleted
)

// Contact 表示系统中的联系人对象，包括用户ID、好友ID、分组ID、备注、状态、置顶标志、黑名单标志和通知启用标志。
type Contact struct {
	*Model              // 模型基类
	User         *User  `json:"-" gorm:"foreignKey:ID;references:FriendId"` // 用户对象
	UserId       int64  `json:"user_id"`                                    // 用户ID
	FriendId     int64  `json:"friend_id"`                                  // 好友ID
	GroupId      int64  `json:"group_id"`                                   // 分组ID
	Remark       string `json:"remark"`                                     // 备注
	Status       int8   `json:"status"`                                     // 状态：1请求好友, 2已同意好友, 3已拒绝好友, 4已删除好友
	IsTop        int8   `json:"is_top"`                                     // 置顶标志
	IsBlack      int8   `json:"is_black"`                                   // 黑名单标志
	NoticeEnable int8   `json:"notice_enable"`                              // 通知启用标志
}

// FetchUser 方法从数据库中获取指定用户和好友的联系记录。
func (c *Contact) FetchUser(db *gorm.DB) (*Contact, error) {
	var contact Contact
	err := db.Omit("User").Unscoped().Where("user_id = ? AND friend_id = ?", c.UserId, c.FriendId).First(&contact).Error
	if err != nil {
		logrus.Debugf("Contact.FetchUser fetch user error:%s", err)
		return nil, err
	}
	return &contact, nil
}

// GetByUserFriend 方法从数据库中获取指定用户和好友的联系记录。
func (c *Contact) GetByUserFriend(db *gorm.DB) (*Contact, error) {
	var contact Contact
	err := db.Omit("User").Where("user_id = ? AND friend_id = ?", c.UserId, c.FriendId).First(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// FetchByUserFriendAll 方法从数据库中获取指定用户和好友的所有联系记录。
func (c *Contact) FetchByUserFriendAll(db *gorm.DB) ([]*Contact, error) {
	var contacts []*Contact
	if err := db.Omit("User").
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			c.UserId, c.FriendId, c.FriendId, c.UserId).
		Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

// List 方法从数据库中按条件获取联系人列表。
func (c *Contact) List(db *gorm.DB, conditions ConditionsT, offset, limit int) ([]*Contact, error) {
	var contacts []*Contact
	var err error
	tn := db.NamingStrategy.TableName("Contact") + "."

	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	for k, v := range conditions {
		if k != "ORDER" {
			db = db.Where(tn+k, v)
		}
	}

	db.Joins("User").Order(clause.OrderByColumn{Column: clause.Column{Name: "nickname"}, Desc: false})
	if err = db.Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

// BeFriendIds 方法获取已同意成为好友关系的好友ID列表。
func (c *Contact) BeFriendIds(db *gorm.DB) (ids []int64, err error) {
	if err = db.Model(c).Omit("User").Select("user_id").Where("friend_id = ? AND status = ?", c.FriendId, ContactStatusAgree).Find(&ids).Error; err != nil {
		return nil, err
	}
	return
}

// MyFriendIds 方法获取当前用户已同意成为好友关系的好友ID列表。
func (c *Contact) MyFriendIds(db *gorm.DB) (ids []string, err error) {
	if err = db.Model(c).Omit("User").Select("friend_id").Where("user_id = ? AND status = ?", c.UserId, ContactStatusAgree).Find(&ids).Error; err != nil {
		return nil, err
	}
	return
}

// Count 方法统计符合条件的联系人数量。
func (m *Contact) Count(db *gorm.DB, conditions ConditionsT) (int64, error) {
	var count int64

	for k, v := range conditions {
		if k != "ORDER" {
			db = db.Where(k, v)
		}
	}
	if err := db.Model(m).Omit("User").Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Create 方法在数据库中创建新的联系人记录。
func (c *Contact) Create(db *gorm.DB) (*Contact, error) {
	err := db.Omit("User").Create(&c).Error
	return c, err
}

// Update 方法更新数据库中的联系人记录。
func (c *Contact) Update(db *gorm.DB) error {
	return db.Model(&Contact{}).Omit("User").Where("id = ?", c.Model.ID).Save(c).Error
}

// UpdateInUnscoped 方法在不考虑软删除的情况下更新数据库中的联系人记录。
func (c *Contact) UpdateInUnscoped(db *gorm.DB) error {
	return db.Unscoped().Omit("User").Save(c).Error
}

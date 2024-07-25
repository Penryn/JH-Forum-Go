// Copyright 2022 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dbr

import (
	"JH-Forum/internal/core/cs"
	"gorm.io/gorm"
)

// 用户状态常量
const (
	UserStatusNormal int = iota + 1
	UserStatusClosed
)

// User 表示用户数据结构
type User struct {
	*Model
	Nickname string `json:"nickname"` // 昵称
	Username string `json:"username"` // 用户名
	Phone    string `json:"phone"`    // 电话
	Password string `json:"password"` // 密码
	Salt     string `json:"salt"`     // 盐
	Status   int    `json:"status"`   // 状态
	Avatar   string `json:"avatar"`   // 头像
	Balance  int64  `json:"balance"`  // 余额
	IsAdmin  bool   `json:"is_admin"` // 是否为管理员
}

// UserFormated 表示格式化后的用户结构
type UserFormated struct {
	ID          int64  `db:"id" json:"id"`   // 用户ID
	Nickname    string `json:"nickname"`     // 昵称
	Username    string `json:"username"`     // 用户名
	Status      int    `json:"status"`       // 状态
	Avatar      string `json:"avatar"`       // 头像
	IsAdmin     bool   `json:"is_admin"`     // 是否为管理员
	IsFriend    bool   `json:"is_friend"`    // 是否为好友
	IsFollowing bool   `json:"is_following"` // 是否正在关注
}

// Format 将用户对象格式化为用户格式化对象
func (u *User) Format() *UserFormated {
	if u.Model != nil {
		return &UserFormated{
			ID:       u.ID,
			Nickname: u.Nickname,
			Username: u.Username,
			Status:   u.Status,
			Avatar:   u.Avatar,
			IsAdmin:  u.IsAdmin,
		}
	}

	return nil
}

// Get 根据条件获取单个用户
func (u *User) Get(db *gorm.DB) (*User, error) {
	var user User
	if u.Model != nil && u.Model.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", u.Model.ID, 0)
	} else if u.Phone != "" {
		db = db.Where("phone = ? AND is_del = ?", u.Phone, 0)
	} else {
		db = db.Where("username = ? AND is_del = ?", u.Username, 0)
	}

	err := db.First(&user).Error
	if err != nil {
		return &user, err
	}

	return &user, nil
}

// List 根据条件获取用户列表
func (u *User) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*User, error) {
	var users []*User
	var err error
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}

	if err = db.Where("is_del = ?", 0).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// ListUserInfoById 根据用户ID列表获取用户信息列表
func (u *User) ListUserInfoById(db *gorm.DB, ids []int64) (res cs.UserInfoList, err error) {
	err = db.Model(u).Where("id IN ?", ids).Find(&res).Error
	return
}

// Create 创建用户
func (u *User) Create(db *gorm.DB) (*User, error) {
	err := db.Create(&u).Error
	return u, err
}

// Update 更新用户信息
func (u *User) Update(db *gorm.DB) error {
	return db.Model(&User{}).Where("id = ? AND is_del = ?", u.Model.ID, 0).Save(u).Error
}

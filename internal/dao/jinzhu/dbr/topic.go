// Copyright 2022 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dbr

import (
	"time"

	"gorm.io/gorm"
)

// Tag 表示标签数据结构。
type Tag struct {
	*Model
	UserID   int64  `json:"user_id"`   // 用户ID
	Tag      string `json:"tag"`       // 标签名
	QuoteNum int64  `json:"quote_num"` // 引用数量
}

type TopicUser struct {
	*Model
	UserID    int64  `json:"user_id"`
	TopicID   int64  `json:"topic_id"`
	AliasName string `json:"-"`
	Remark    string `json:"-"`
	QuoteNum  int64  `json:"quote_num"`
	IsTop     int8   `json:"is_top"`
	ReserveA  string `json:"-"`
	ReserveB  string `json:"-"`
}

// TagFormated 表示格式化后的标签结构。
type TagFormated struct {
	ID          int64         `json:"id"`           // 标签ID
	UserID      int64         `json:"user_id"`      // 用户ID
	User        *UserFormated `json:"user"`         // 用户信息
	Tag         string        `json:"tag"`          // 标签名
	QuoteNum    int64         `json:"quote_num"`    // 引用数量
	IsFollowing int8          `json:"is_following"` // 是否正在关注
	IsTop       int8          `json:"is_top"`       // 是否置顶
}

// Format 将标签对象格式化为标签格式化对象。
func (t *Tag) Format() *TagFormated {
	if t.Model == nil {
		return &TagFormated{}
	}

	return &TagFormated{
		ID:          t.ID,
		UserID:      t.UserID,
		User:        &UserFormated{},
		Tag:         t.Tag,
		QuoteNum:    t.QuoteNum,
		IsFollowing: 0,
		IsTop:       0,
	}
}

// Get 根据条件获取单个标签。
func (t *Tag) Get(db *gorm.DB) (*Tag, error) {
	var tag Tag
	if t.Model != nil && t.Model.ID > 0 {
		db = db.Where("id = ? AND is_del = ?", t.Model.ID, 0)
	} else {
		db = db.Where("tag = ? AND is_del = ?", t.Tag, 0)
	}

	err := db.First(&tag).Error
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// Create 创建标签。
func (t *Tag) Create(db *gorm.DB) (*Tag, error) {
	err := db.Create(&t).Error
	return t, err
}

// Update 更新标签信息。
func (t *Tag) Update(db *gorm.DB) error {
	return db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.Model.ID, 0).Save(t).Error
}

// Delete 根据ID删除标签。
func (t *Tag) Delete(db *gorm.DB) error {
	return db.Model(t).Where("id = ?", t.Model.ID).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// List 根据条件获取标签列表。
func (t *Tag) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) (tags []*Tag, err error) {
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if t.UserID > 0 {
		db = db.Where("user_id = ?", t.UserID)
	}
	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}
	err = db.Where("is_del = 0 and quote_num > 0").Find(&tags).Error
	return
}

// TagsFrom 根据标签名数组获取标签列表。
func (t *Tag) TagsFrom(db *gorm.DB, tags []string) (res []*Tag, err error) {
	err = db.Where("tag IN ?", tags).Find(&res).Error
	return
}

// PostCollection 定义了帖子收藏模型，包含帖子ID、用户ID等信息。
package dbr

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PostCollection 表示帖子收藏的信息。
type PostCollection struct {
	*Model
	Post   *Post `json:"-"`                    // 关联的帖子信息
	PostID int64 `db:"post_id" json:"post_id"` // 帖子ID
	UserID int64 `db:"user_id" json:"user_id"` // 用户ID
}

// Get 根据条件获取单个帖子收藏记录。
func (p *PostCollection) Get(db *gorm.DB) (*PostCollection, error) {
	var star PostCollection
	tn := db.NamingStrategy.TableName("PostCollection") + "."

	if p.Model != nil && p.ID > 0 {
		db = db.Where(tn+"id = ? AND "+tn+"is_del = ?", p.ID, 0)
	}
	if p.PostID > 0 {
		db = db.Where(tn+"post_id = ?", p.PostID)
	}
	if p.UserID > 0 {
		db = db.Where(tn+"user_id = ?", p.UserID)
	}

	// 根据帖子的可见性进行过滤和排序
	db = db.Joins("Post").Where("visibility <> ? OR (visibility = ? AND ? = ?)", PostVisitPrivate, PostVisitPrivate, clause.Column{Table: "Post", Name: "user_id"}, p.UserID).Order(clause.OrderByColumn{Column: clause.Column{Table: "Post", Name: "id"}, Desc: true})
	err := db.First(&star).Error
	if err != nil {
		return &star, err
	}

	return &star, nil
}

// Create 创建帖子收藏记录。
func (p *PostCollection) Create(db *gorm.DB) (*PostCollection, error) {
	err := db.Omit("Post").Create(&p).Error
	return p, err
}

// Delete 根据ID删除帖子收藏记录。
func (p *PostCollection) Delete(db *gorm.DB) error {
	return db.Model(&PostCollection{}).Omit("Post").Where("id = ? AND is_del = ?", p.Model.ID, 0).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// List 根据条件查询帖子收藏记录列表。
func (p *PostCollection) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*PostCollection, error) {
	var collections []*PostCollection
	var err error
	tn := db.NamingStrategy.TableName("PostCollection") + "."

	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if p.UserID > 0 {
		db = db.Where(tn+"user_id = ?", p.UserID)
	}

	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(tn+k, v)
		}
	}

	// 根据帖子的可见性进行过滤和排序
	db = db.Joins("Post").Where(`visibility <> ? OR (visibility = ? AND ? = ?)`, PostVisitPrivate, PostVisitPrivate, clause.Column{Table: "Post", Name: "user_id"}, p.UserID).Order(clause.OrderByColumn{Column: clause.Column{Table: "Post", Name: "id"}, Desc: true})
	if err = db.Where(tn+"is_del = ?", 0).Find(&collections).Error; err != nil {
		return nil, err
	}

	return collections, nil
}

// Count 根据条件统计帖子收藏记录数量。
func (p *PostCollection) Count(db *gorm.DB, conditions *ConditionsT) (int64, error) {
	var count int64
	tn := db.NamingStrategy.TableName("PostCollection") + "."

	if p.PostID > 0 {
		db = db.Where(tn+"post_id = ?", p.PostID)
	}
	if p.UserID > 0 {
		db = db.Where(tn+"user_id = ?", p.UserID)
	}
	for k, v := range *conditions {
		if k != "ORDER" {
			db = db.Where(tn+k, v)
		}
	}

	// 根据帖子的可见性进行过滤
	db = db.Joins("Post").Where(`visibility <> ? OR (visibility = ? AND ? = ?)`, PostVisitPrivate, PostVisitPrivate, clause.Column{Table: "Post", Name: "user_id"}, p.UserID)
	if err := db.Model(p).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

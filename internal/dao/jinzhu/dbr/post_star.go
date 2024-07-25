// PostStar 定义了帖子点赞模型，包含帖子ID、用户ID等信息。
package dbr

import (
	"time"

	"JH-Forum/internal/core/cs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PostStar 表示帖子点赞记录。
type PostStar struct {
	*Model
	Post   *Post `json:"-"`
	PostID int64 `json:"post_id"` // 帖子ID
	UserID int64 `json:"user_id"` // 用户ID
}

// Get 根据条件获取单个帖子点赞记录。
func (p *PostStar) Get(db *gorm.DB) (*PostStar, error) {
	var star PostStar
	tn := db.NamingStrategy.TableName("PostStar") + "."

	if p.Model != nil && p.ID > 0 {
		db = db.Where(tn+"id = ? AND "+tn+"is_del = ?", p.ID, 0)
	}
	if p.PostID > 0 {
		db = db.Where(tn+"post_id = ?", p.PostID)
	}
	if p.UserID > 0 {
		db = db.Where(tn+"user_id = ?", p.UserID)
	}

	db = db.Joins("Post").Where("visibility <> ? OR (visibility = ? AND ? = ?)", PostVisitPrivate, PostVisitPrivate, clause.Column{Table: "Post", Name: "user_id"}, p.UserID).Order(clause.OrderByColumn{Column: clause.Column{Table: "Post", Name: "id"}, Desc: true})
	if err := db.First(&star).Error; err != nil {
		return nil, err
	}
	return &star, nil
}

// Create 创建帖子点赞记录。
func (p *PostStar) Create(db *gorm.DB) (*PostStar, error) {
	err := db.Omit("Post").Create(&p).Error
	return p, err
}

// Delete 删除帖子点赞记录。
func (p *PostStar) Delete(db *gorm.DB) error {
	return db.Model(&PostStar{}).Omit("Post").Where("id = ? AND is_del = ?", p.Model.ID, 0).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}

// List 根据条件查询帖子点赞记录列表。
func (p *PostStar) List(db *gorm.DB, conditions *ConditionsT, typ cs.RelationTyp, limit int, offset int) (res []*PostStar, err error) {
	tn := db.NamingStrategy.TableName("PostStar") + "."
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
	db = db.Joins("Post")
	switch typ {
	case cs.RelationAdmin:
		// 管理员具有访问所有类型帖子的权限
	case cs.RelationFriend:
		db = db.Where("visibility = ? OR visibility = ?", PostVisitPublic, PostVisitFriend)
	case cs.RelationSelf:
		db = db.Where("visibility <> ? OR (visibility = ? AND ? = ?)", PostVisitPrivate, PostVisitPrivate, clause.Column{Table: "Post", Name: "user_id"}, p.UserID)
	default:
		db = db.Where("visibility=?", PostVisitPublic)
	}
	db = db.Order(clause.OrderByColumn{Column: clause.Column{Table: "Post", Name: "id"}, Desc: true})
	err = db.Find(&res).Error
	return
}

// Count 根据条件统计帖子点赞记录数量。
func (p *PostStar) Count(db *gorm.DB, typ cs.RelationTyp, conditions *ConditionsT) (res int64, err error) {
	tn := db.NamingStrategy.TableName("PostStar") + "."
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
	db = db.Joins("Post")
	switch typ {
	case cs.RelationAdmin:
		// 管理员具有访问所有类型帖子的权限
	case cs.RelationFriend:
		db = db.Where("visibility = ? OR visibility = ?", PostVisitPublic, PostVisitFriend)
	case cs.RelationSelf:
		db = db.Where("visibility <> ? OR (visibility = ? AND ? = ?)", PostVisitPrivate, PostVisitPrivate, clause.Column{Table: "Post", Name: "user_id"}, p.UserID)
	default:
		db = db.Where("visibility=?", PostVisitPublic)
	}
	err = db.Model(p).Count(&res).Error
	return
}

package dbr

import "gorm.io/gorm"

// WalletStatement 表示钱包交易记录
type WalletStatement struct {
	*Model
	UserID          int64  `json:"user_id"`          // 用户ID
	ChangeAmount    int64  `json:"change_amount"`    // 变动金额
	BalanceSnapshot int64  `json:"balance_snapshot"` // 余额快照
	Reason          string `json:"reason"`           // 变动原因
	PostID          int64  `json:"post_id"`          // 相关的帖子ID
}

// Get 根据条件获取单个钱包交易记录
func (w *WalletStatement) Get(db *gorm.DB) (*WalletStatement, error) {
	var statement WalletStatement
	query := db.Where("is_del = ?", 0)
	if w.Model != nil && w.ID > 0 {
		query = query.Where("id = ?", w.ID)
	}
	if w.PostID > 0 {
		query = query.Where("post_id = ?", w.PostID)
	}
	if w.UserID > 0 {
		query = query.Where("user_id = ?", w.UserID)
	}

	err := query.First(&statement).Error
	if err != nil {
		return &statement, err
	}

	return &statement, nil
}

// Create 创建钱包交易记录
func (w *WalletStatement) Create(db *gorm.DB) (*WalletStatement, error) {
	err := db.Create(&w).Error
	return w, err
}

// List 获取钱包交易记录列表
func (w *WalletStatement) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*WalletStatement, error) {
	var statements []*WalletStatement
	query := db.Where("is_del = ?", 0)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	if w.UserID > 0 {
		query = query.Where("user_id = ?", w.UserID)
	}

	for k, v := range *conditions {
		if k == "ORDER" {
			query = query.Order(v)
		} else {
			query = query.Where(k, v)
		}
	}

	err := query.Find(&statements).Error
	if err != nil {
		return nil, err
	}

	return statements, nil
}

// Count 获取钱包交易记录的数量
func (w *WalletStatement) Count(db *gorm.DB, conditions *ConditionsT) (int64, error) {
	var count int64
	query := db.Model(w).Where("is_del = ?", 0)
	if w.PostID > 0 {
		query = query.Where("post_id = ?", w.PostID)
	}
	if w.UserID > 0 {
		query = query.Where("user_id = ?", w.UserID)
	}

	for k, v := range *conditions {
		if k != "ORDER" {
			query = query.Where(k, v)
		}
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

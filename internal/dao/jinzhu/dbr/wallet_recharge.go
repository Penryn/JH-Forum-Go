package dbr

import "gorm.io/gorm"

// WalletRecharge 表示钱包充值记录
type WalletRecharge struct {
	*Model
	UserID      int64  `json:"user_id"`      // 用户ID
	Amount      int64  `json:"amount"`       // 充值金额
	TradeNo     string `json:"trade_no"`     // 交易号
	TradeStatus string `json:"trade_status"` // 交易状态
}

// Get 根据条件获取单个钱包充值记录
func (p *WalletRecharge) Get(db *gorm.DB) (*WalletRecharge, error) {
	var recharge WalletRecharge
	query := db.Where("is_del = ?", 0)
	if p.Model != nil && p.ID > 0 {
		query = query.Where("id = ?", p.ID)
	}
	if p.UserID > 0 {
		query = query.Where("user_id = ?", p.UserID)
	}

	err := query.First(&recharge).Error
	if err != nil {
		return &recharge, err
	}

	return &recharge, nil
}

// Create 创建钱包充值记录
func (p *WalletRecharge) Create(db *gorm.DB) (*WalletRecharge, error) {
	err := db.Create(&p).Error
	return p, err
}

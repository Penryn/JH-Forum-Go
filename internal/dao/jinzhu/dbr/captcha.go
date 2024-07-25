package dbr

import "gorm.io/gorm"

// Captcha 结构体表示验证码对象，包括手机号、验证码内容、使用次数和过期时间。
type Captcha struct {
	*Model           // 模型基类
	Phone     string `json:"phone"`      // 手机号
	Captcha   string `json:"captcha"`    // 验证码
	UseTimes  int    `json:"use_times"`  // 使用次数
	ExpiredOn int64  `json:"expired_on"` // 过期时间
}

// Create 在数据库中创建验证码记录。
// 返回创建的验证码对象和可能发生的错误。
func (c *Captcha) Create(db *gorm.DB) (*Captcha, error) {
	err := db.Create(&c).Error
	return c, err
}

// Update 更新数据库中的验证码记录。
// 返回可能发生的错误。
func (c *Captcha) Update(db *gorm.DB) error {
	return db.Model(&Captcha{}).Where("id = ? AND is_del = ?", c.Model.ID, 0).Save(c).Error
}

// Get 根据条件从数据库中获取验证码记录。
// 如果设置了 ID，则根据 ID 和 is_del=0 进行查询；
// 如果设置了 Phone，则根据 Phone 进行查询。
// 返回获取到的验证码对象和可能发生的错误。
func (c *Captcha) Get(db *gorm.DB) (*Captcha, error) {
	var captcha Captcha
	if c.Model != nil && c.ID > 0 {
		db = db.Where("id = ? AND is_del = ?", c.ID, 0)
	}
	if c.Phone != "" {
		db = db.Where("phone = ?", c.Phone)
	}

	err := db.Last(&captcha).Error
	if err != nil {
		return &captcha, err
	}

	return &captcha, nil
}

// Model 定义了一个公共模型结构，包含标准字段如ID、时间戳和软删除状态。
package dbr

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// Model 表示一个通用的模型结构，包含标准字段如ID、时间戳和软删除。
type Model struct {
	ID         int64                 `gorm:"primary_key" json:"id"`         // 主键ID
	CreatedOn  int64                 `json:"created_on"`                    // 创建时间戳
	ModifiedOn int64                 `json:"modified_on"`                   // 修改时间戳
	DeletedOn  int64                 `json:"deleted_on"`                    // 删除时间戳
	IsDel      soft_delete.DeletedAt `gorm:"softDelete:flag" json:"is_del"` // 软删除标志
}

// ConditionsT 定义数据库查询条件类型。
type ConditionsT map[string]interface{}

// Predicates 定义数据库查询条件谓词。
type Predicates map[string][]interface{}

// BeforeCreate 在创建记录之前设置创建时间和修改时间戳。
func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	nowTime := time.Now().Unix()

	tx.Statement.SetColumn("created_on", nowTime)
	tx.Statement.SetColumn("modified_on", nowTime)
	return
}

// BeforeUpdate 在更新记录之前设置修改时间戳。
func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	if !tx.Statement.Changed("modified_on") {
		tx.Statement.SetColumn("modified_on", time.Now().Unix())
	}

	return
}

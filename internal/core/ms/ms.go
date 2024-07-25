// 该代码定义了一些常量和类型的别名，与外部包中的类型对应。

// Package ms 包含了核心数据服务接口类型定义，为gorm适配器定义模型。
package ms

import (
	"JH-Forum/internal/dao/jinzhu/dbr"
)

const (
	UserStatusNormal = dbr.UserStatusNormal // UserStatusNormal 的别名
	UserStatusClosed = dbr.UserStatusClosed // UserStatusClosed 的别名
)

type (
	User                = dbr.User                // User 的别名
	Post                = dbr.Post                // Post 的别名
	ConditionsT         = dbr.ConditionsT         // ConditionsT 的别名
	PostFormated        = dbr.PostFormated        // PostFormated 的别名
	UserFormated        = dbr.UserFormated        // UserFormated 的别名
	PostContentFormated = dbr.PostContentFormated // PostContentFormated 的别名
	Model               = dbr.Model               // Model 的别名
)

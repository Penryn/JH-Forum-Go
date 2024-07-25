// 该代码定义了一些内部核心错误变量，用于数据逻辑实现。

package cs

import "errors"

// ErrNotImplemented 表示未实现的错误
var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNoPermission   = errors.New("no permission")
)

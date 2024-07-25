// 该代码定义了一个用于版本信息的接口。

package core

import (
	"github.com/Masterminds/semver/v3"
)

// VersionInfo 版本信息接口
type VersionInfo interface {
	Name() string             // 获取版本名称
	Version() *semver.Version // 获取版本号
}

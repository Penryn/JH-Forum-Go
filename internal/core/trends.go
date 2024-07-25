// 该代码定义了一个接口，用于实现动态信息管理服务。

package core

import (
	"JH-Forum/internal/core/cs"
)

// TrendsManageServantA 动态信息管理服务接口
type TrendsManageServantA interface {
	GetIndexTrends(userId int64, limit int, offset int) ([]*cs.TrendsItem, int64, error) // 获取首页动态信息
}

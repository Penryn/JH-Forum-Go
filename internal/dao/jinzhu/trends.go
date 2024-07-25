// package jinzhu 实现了广场推文索引服务，包括根据用户ID查询广场推文列表等功能。

package jinzhu

import (
	"fmt" // 格式化输出包

	"JH-Forum/internal/core"    // 引入核心服务包
	"JH-Forum/internal/core/cs" // 引入核心服务-通用包
	"gorm.io/gorm"              // 引入gorm包
)

// trendsSrvA 提供处理趋势管理服务的实现。
type trendsSrvA struct {
	db *gorm.DB // GORM数据库连接
}

// GetIndexTrends 获取指定用户的首页趋势数据列表。
// 根据用户ID查询相关的趋势信息，包括用户名称、昵称和头像等。
func (s *trendsSrvA) GetIndexTrends(userId int64, limit int, offset int) (res []*cs.TrendsItem, total int64, err error) {
	db := s.db.Table(_user_).
		Joins(fmt.Sprintf("JOIN %s r ON r.he_uid=%s.id", _userRelation_, _user_)).
		Joins(fmt.Sprintf("JOIN %s m ON r.he_uid=m.user_id", _userMetric_)).
		Where("r.user_id=? AND m.tweets_count>0 AND m.is_del=0", userId)

	if err = db.Count(&total).Error; err != nil || total == 0 {
		return
	}

	if offset >= 0 && limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Order("r.style ASC, m.latest_trends_on DESC").Select("username", "nickname", "avatar").Find(&res).Error; err == nil {
		res = cs.DistinctTrends(res)
	}

	return
}

// newTrendsManageServentA 创建一个新的趋势管理服务实例。
// 返回一个可以操作趋势管理服务的对象。
func newTrendsManageServentA(db *gorm.DB) core.TrendsManageServantA {
	return &trendsSrvA{
		db: db,
	}
}

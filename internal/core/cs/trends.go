// 该代码定义了用于处理趋势项的结构体和函数。

package cs

import "JH-Forum/pkg/types"

// TrendsItem 表示趋势项的结构体
type TrendsItem struct {
	Username string `json:"username"`          // 用户名
	Nickname string `json:"nickname"`          // 昵称
	Avatar   string `json:"avatar"`            // 头像链接
	IsFresh  bool   `json:"is_fresh" gorm:"-"` // 是否新鲜，不作为数据库字段
}

// DistinctTrends 根据用户名去重趋势项
func DistinctTrends(items []*TrendsItem) []*TrendsItem {
	if len(items) == 0 {
		return items
	}
	res := make([]*TrendsItem, 0, len(items))
	set := make(map[string]types.Empty, len(items))
	for _, item := range items {
		if _, exist := set[item.Username]; exist {
			continue
		}
		res = append(res, item)
		set[item.Username] = types.Empty{}
	}
	return res
}

// 该代码定义了标签相关的常量、类型和结构体。

package cs

const (
	// 标签类型
	TagTypeHot       TagType = "hot"        // 热门标签类型
	TagTypeNew       TagType = "new"        // 新标签类型
	TagTypeFollow    TagType = "follow"     // 关注标签类型
	TagTypeHotExtral TagType = "hot_extral" // 额外热门标签类型
)

type (
	// TagType 标签类型
	TagType string

	// TagInfoList 标签信息列表
	TagInfoList []*TagInfo

	// TagList 标签列表
	TagList []*TagItem
)

// TagInfo 表示标签信息的结构体
type TagInfo struct {
	ID       int64  `json:"id"`        // 标签ID
	UserID   int64  `json:"user_id"`   // 用户ID
	Tag      string `json:"tag"`       // 标签名称
	QuoteNum int64  `json:"quote_num"` // 引用数
}

// TagItem 表示标签信息条目的结构体
type TagItem struct {
	ID          int64     `json:"id"`           // 标签ID
	UserID      int64     `json:"user_id"`      // 用户ID
	User        *UserInfo `json:"user"`         // 用户信息
	Tag         string    `json:"tag"`          // 标签名称
	QuoteNum    int64     `json:"quote_num"`    // 引用数
	IsFollowing int8      `json:"is_following"` // 是否关注
	IsTop       int8      `json:"is_top"`       // 是否置顶
}

// Format 将 TagInfo 格式化为 TagItem
func (t *TagInfo) Format() *TagItem {
	return &TagItem{
		ID:          t.ID,
		UserID:      t.UserID,
		User:        &UserInfo{},
		Tag:         t.Tag,
		QuoteNum:    t.QuoteNum,
		IsFollowing: 0,
		IsTop:       0,
	}
}

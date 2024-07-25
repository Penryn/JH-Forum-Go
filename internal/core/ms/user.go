// 该代码定义了 ContactItem 和 ContactList 结构体，用于表示联系人信息和联系人列表信息。

package ms

type (
	// ContactItem 表示单个联系人的信息
	ContactItem struct {
		UserId      int64  `json:"user_id"`                // 用户ID
		Username    string `db:"username" json:"username"` // 用户名
		Nickname    string `json:"nickname"`               // 昵称
		Avatar      string `json:"avatar"`                 // 头像地址
		Phone       string `json:"phone,omitempty"`        // 手机号（可选）
		IsFollowing bool   `json:"is_following"`           // 是否正在关注该联系人
		CreatedOn   int64  `json:"created_on"`             // 创建时间戳
	}

	// ContactList 表示联系人列表信息
	ContactList struct {
		Contacts []ContactItem `json:"contacts"` // 联系人列表
		Total    int64         `json:"total"`    // 总数
	}
)

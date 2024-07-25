// 该代码定义了用户关系、用户信息和用户个人资料相关的常量、类型和结构体。

package cs

const (
	RelationUnknown   RelationTyp = iota // 0 未知关系
	RelationSelf                         // 1 自己
	RelationFriend                       // 2 好友
	RelationFollower                     // 3 粉丝
	RelationFollowing                    // 4 关注的人
	RelationAdmin                        // 5 管理员
	RelationGuest                        // 6 游客
)

type (
	// UserInfoList 用户信息列表
	UserInfoList []*UserInfo

	// RelationTyp 表示用户关系类型
	RelationTyp uint8

	// VistUser 表示访问用户信息
	VistUser struct {
		Username string      // 用户名
		UserId   int64       // 用户ID
		RelTyp   RelationTyp // 关系类型
	}
)

// UserInfo 表示用户基本信息的结构体
type UserInfo struct {
	ID        int64  `json:"id"`         // 用户ID
	Nickname  string `json:"nickname"`   // 昵称
	Username  string `json:"username"`   // 用户名
	Status    int    `json:"status"`     // 状态
	Avatar    string `json:"avatar"`     // 头像
	IsAdmin   bool   `json:"is_admin"`   // 是否管理员
	CreatedOn int64  `json:"created_on"` // 创建时间
}

// UserProfile 表示用户个人资料的结构体
type UserProfile struct {
	ID          int64  `json:"id"`           // 用户ID
	Nickname    string `json:"nickname"`     // 昵称
	Username    string `json:"username"`     // 用户名
	Phone       string `json:"phone"`        // 手机号
	Status      int    `json:"status"`       // 状态
	Avatar      string `json:"avatar"`       // 头像
	Balance     int64  `json:"balance"`      // 余额
	IsAdmin     bool   `json:"is_admin"`     // 是否管理员
	CreatedOn   int64  `json:"created_on"`   // 创建时间
	TweetsCount int    `json:"tweets_count"` // 推文数
}

// String 返回 RelationTyp 的字符串表示
func (t RelationTyp) String() string {
	switch t {
	case RelationSelf:
		return "self"
	case RelationFriend:
		return "friend"
	case RelationFollower:
		return "follower"
	case RelationFollowing:
		return "following"
	case RelationAdmin:
		return "admin"
	case RelationUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

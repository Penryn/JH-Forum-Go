// 该代码定义了联系人状态常量和联系人结构体。

package cs

const (
	// ContactStatusRequesting 表示请求添加好友状态
	ContactStatusRequesting int8 = iota + 1
	// ContactStatusAgree 表示已同意好友状态
	ContactStatusAgree
	// ContactStatusReject 表示已拒绝好友状态
	ContactStatusReject
	// ContactStatusDeleted 表示已删除好友状态
	ContactStatusDeleted
)

// Contact 表示联系人结构体
type Contact struct {
	ID           int64  `db:"id" json:"id"`               // ID
	UserId       int64  `db:"user_id" json:"user_id"`     // 用户ID
	FriendId     int64  `db:"friend_id" json:"friend_id"` // 好友ID
	GroupId      int64  `json:"group_id"`                 // 分组ID
	Remark       string `json:"remark"`                   // 备注
	Status       int8   `json:"status"`                   // 状态：1请求好友, 2已同意好友, 3已拒绝好友, 4已删除好友
	IsTop        int8   `json:"is_top"`                   // 是否置顶
	IsBlack      int8   `json:"is_black"`                 // 是否加入黑名单
	NoticeEnable int8   `json:"notice_enable"`            // 是否开启通知
	IsDel        int8   `json:"-"`                        // 是否删除
	DeletedOn    int64  `db:"-" json:"-"`                 // 删除时间
}

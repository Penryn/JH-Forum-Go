// 该代码定义了附件账单的结构体。

package cs

// AttachmentBill 表示附件账单
type AttachmentBill struct {
	ID         int64 `json:"id"`          // 账单ID
	PostID     int64 `json:"post_id"`     // 所属帖子ID
	UserID     int64 `json:"user_id"`     // 用户ID
	PaidAmount int64 `json:"paid_amount"` // 支付金额
}

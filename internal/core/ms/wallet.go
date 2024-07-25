// 该代码导入了 "github.com/rocboss/paopao-ce/internal/dao/jinzhu/dbr" 包，并定义了 WalletStatement 和 WalletRecharge 结构体。

package ms

import (
	"JH-Forum/internal/dao/jinzhu/dbr"
)

type (
	// WalletStatement 表示钱包账单信息
	WalletStatement = dbr.WalletStatement

	// WalletRecharge 表示钱包充值信息
	WalletRecharge = dbr.WalletRecharge
)

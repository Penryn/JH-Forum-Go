// Copyright 2023 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chain

import (
	"sync"

	"JH-Forum/internal/core"
	"JH-Forum/internal/dao"
	"JH-Forum/internal/dao/cache"
)

var (
	_ums     core.UserManageService
	_ac      core.AppCache
	_onceUms sync.Once
)

func userManageService() core.UserManageService {
	_onceUms.Do(func() {
		_ums = dao.DataService()
		_ac = cache.NewAppCache()
	})
	return _ums
}

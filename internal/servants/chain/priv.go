// Copyright 2022 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chain

import (

	"github.com/gin-gonic/gin"
	"JH-Forum/internal/core/ms"
	"JH-Forum/pkg/app"
)

func Priv() gin.HandlerFunc {
	return func(c *gin.Context) {
		if u, exist := c.Get("USER"); exist {
			if user, ok := u.(*ms.User); ok && user.Status == ms.UserStatusNormal {
				c.Next()
				return
			}
		}
		response := app.NewResponse(c)
		response.ToErrorResponse(_errUserHasBeenBanned)
		c.Abort()
	}
}

// Copyright 2022 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"JH-Forum/internal/core/ms"
)

// SecurityService 安全相关服务
type SecurityService interface {
	GetLatestPhoneCaptcha(phone string) (*ms.Captcha, error)
}

// AttachmentCheckService 附件检测服务
type AttachmentCheckService interface {
	CheckAttachment(uri string) error
}

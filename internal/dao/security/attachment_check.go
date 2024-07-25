// 包 security 实现了与安全相关的服务功能。

package security

import (
	"fmt"
	"strings"

	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
)

// attachmentCheckServant 实现了附件检查服务。
type attachmentCheckServant struct {
	domain string
}

// CheckAttachment 检查附件是否为本站资源。
func (s *attachmentCheckServant) CheckAttachment(uri string) error {
	if !strings.HasPrefix(uri, s.domain) {
		return fmt.Errorf("附件非本站资源")
	}
	return nil
}

// NewAttachmentCheckService 创建一个新的附件检查服务。
func NewAttachmentCheckService() core.AttachmentCheckService {
	return &attachmentCheckServant{
		domain: conf.GetOssDomain(),
	}
}

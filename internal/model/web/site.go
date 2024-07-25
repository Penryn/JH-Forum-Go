package web

import (
	"JH-Forum/internal/conf"
	"JH-Forum/pkg/version"
)

// VersionResp 定义了版本信息响应结构体。
type VersionResp struct {
	BuildInfo *version.BuildInfo `json:"build_info"`
}

// SiteProfileResp 是站点配置响应类型的别名，使用 conf.WebProfileConf 类型。
type SiteProfileResp = conf.WebProfileConf

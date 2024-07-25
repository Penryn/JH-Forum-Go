// Copyright 2023 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package web

import (
	"github.com/alimy/mir/v4"
	api "JH-Forum/mirc/auto/api/v1"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/model/web"
	"JH-Forum/internal/servants/base"
	"JH-Forum/pkg/version"
)

var (
	_ api.Site = (*siteSrv)(nil)
)

type siteSrv struct {
	api.UnimplementedSiteServant
	*base.BaseServant
}

func (*siteSrv) Profile() (*web.SiteProfileResp, mir.Error) {
	return conf.WebProfileSetting, nil
}

func (*siteSrv) Version() (*web.VersionResp, mir.Error) {
	return &web.VersionResp{
		BuildInfo: version.ReadBuildInfo(),
	}, nil
}

func newSiteSrv() api.Site {
	return &siteSrv{
		BaseServant: base.NewBaseServant(),
	}
}

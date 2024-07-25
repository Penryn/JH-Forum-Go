package dao

import (
	"sync"

	"github.com/alimy/tryst/cfg"
	"JH-Forum/internal/core"
	"JH-Forum/internal/dao/jinzhu"
	"JH-Forum/internal/dao/search"

	"JH-Forum/internal/dao/storage"
	"github.com/sirupsen/logrus"
)

var (
	ts     core.TweetSearchService
	ds     core.DataService
	oss    core.ObjectStorageService
	webDsa core.WebDataServantA

	_onceInitial sync.Once
)

// DataService 返回核心数据服务
func DataService() core.DataService {
	lazyInitial()
	return ds
}

// WebDataServantA 返回Web数据服务A
func WebDataServantA() core.WebDataServantA {
	lazyInitial()
	return webDsa
}

// ObjectStorageService 返回对象存储服务
func ObjectStorageService() core.ObjectStorageService {
	lazyInitial()
	return oss
}

// TweetSearchService 返回推特搜索服务
func TweetSearchService() core.TweetSearchService {
	lazyInitial()
	return ts
}

// newAuthorizationManageService 创建授权管理服务实例
func newAuthorizationManageService() core.AuthorizationManageService {
	ams := jinzhu.NewAuthorizationManageService()
	return ams
}

// lazyInitial 执行一些包的延迟初始化以提升性能
func lazyInitial() {
	_onceInitial.Do(func() {
		initDsX()
		initOSS()
		initTsX()
	})
}

// initDsX 初始化数据服务和Web数据服务A
func initDsX() {
	var dsVer, dsaVer core.VersionInfo
	ds, dsVer = jinzhu.NewDataService()
	webDsa, dsaVer = jinzhu.NewWebDataServantA()
	logrus.Infof("使用 %s 作为核心数据服务，版本 %s", dsVer.Name(), dsVer.Version())
	logrus.Infof("使用 %s 作为Web数据服务A，版本 %s", dsaVer.Name(), dsaVer.Version())
}

// initOSS 初始化对象存储服务
func initOSS() {
	var v core.VersionInfo
	oss, v = storage.MustMinioService()
	logrus.Infof("使用 %s 作为对象存储服务，版本 %s", v.Name(), v.Version())
}

// initTsX 初始化推文搜索服务
func initTsX() {
	var v core.VersionInfo
	ams := newAuthorizationManageService()
	cfg.On(cfg.Actions{
		"Meili": func() {
			ts, v = search.NewMeiliTweetSearchService(ams)
		},
	}, func() {
		ts, v = search.NewMeiliTweetSearchService(ams)
	})
	logrus.Infof("使用 %s 作为推文搜索服务，版本 %s", v.Name(), v.Version())
	ts = search.NewBridgeTweetSearchService(ts)
}

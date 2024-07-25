package web

import (
	"sync"

	"github.com/alimy/tryst/cfg"
	"github.com/gin-gonic/gin"
	api "JH-Forum/mirc/auto/api/v1"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"JH-Forum/internal/dao"
	"JH-Forum/internal/dao/cache"
	"JH-Forum/internal/servants/base"
)

var (
	_disallowUserRegister bool
	_ds                   core.DataService
	_ac                   core.AppCache
	_wc                   core.WebCache
	_oss                  core.ObjectStorageService
	_onceInitial          sync.Once
)

// RouteWeb 注册 web 路由
func RouteWeb(e *gin.Engine) {
	lazyInitial()
	ds := base.NewDaoServant()
	// 始终注册 servants
	api.RegisterAdminServant(e, newAdminSrv(ds, _wc))
	api.RegisterCoreServant(e, newCoreSrv(ds, _oss, _wc))
	api.RegisterRelaxServant(e, newRelaxSrv(ds, _wc), newRelaxChain())
	api.RegisterLooseServant(e, newLooseSrv(ds, _ac))
	api.RegisterPrivServant(e, newPrivSrv(ds, _oss), newPrivChain())
	api.RegisterPubServant(e, newPubSrv(ds))
	api.RegisterTrendsServant(e, newTrendsSrv(ds))
	api.RegisterFollowshipServant(e, newFollowshipSrv(ds))
	api.RegisterFriendshipServant(e, newFriendshipSrv(ds))
	api.RegisterSiteServant(e, newSiteSrv())
	// 根据配置注册所需的 servants
	// 如需则调度作业
	scheduleJobs()
}

// lazyInitial 执行一些包的延迟初始化以提升性能
func lazyInitial() {
	_onceInitial.Do(func() {
		_disallowUserRegister = cfg.If("Web:DisallowUserRegister")
		_maxWhisperNumDaily = conf.AppSetting.MaxWhisperDaily
		_oss = dao.ObjectStorageService()
		_ds = dao.DataService()
		_ac = cache.NewAppCache()
		_wc = cache.NewWebCache()
	})
}

package prometheus

import (
	"github.com/prometheus/client_golang/prometheus" // 导入Prometheus客户端库
	"JH-Forum/internal/conf"                         // 导入配置包
	"JH-Forum/internal/core"                         // 导入核心包
	"github.com/sirupsen/logrus"                     // 导入日志库
)

// metrics 结构定义了度量指标
type metrics struct {
	siteInfo *prometheus.GaugeVec // 网站信息度量指标向量
	ds       core.DataService     // 数据服务
	wc       core.WebCache        // Web缓存
}

// updateSiteInfo 更新网站信息度量指标
func (m *metrics) updateSiteInfo() {
	// 更新最大在线用户数
	if onlineUserKeys, err := m.wc.Keys(conf.PrefixOnlineUser + "*"); err == nil {
		maxOnline := len(onlineUserKeys)
		m.siteInfo.With(prometheus.Labels{"name": "max_online"}).Set(float64(maxOnline))
	} else {
		logrus.Warnf("update promethues metrics[site_info_max_online] occurs error: %s", err)
	}
	// 更新注册用户数
	if registerUserCount, err := m.ds.GetRegisterUserCount(); err == nil {
		m.siteInfo.With(prometheus.Labels{"name": "register_user_count"}).Set(float64(registerUserCount))
	} else {
		logrus.Warnf("update promethues metrics[site_info_register_user_count] occurs error: %s", err)
	}
}

// onUpdate 执行度量指标更新作业
func (m *metrics) onUpdate() {
	logrus.Debugf("update promethues metrics job running") // 记录调试信息
	m.updateSiteInfo()                                     // 更新网站信息度量指标
}

// newMetrics 创建新的度量指标管理器
func newMetrics(reg prometheus.Registerer, ds core.DataService, wc core.WebCache) *metrics {
	m := &metrics{
		siteInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "paopao",
				Subsystem: "site",
				Name:      "simple_info",
				Help:      "paopao-ce site simple information.", // 网站简要信息
			},
			[]string{
				"name", // 指标名称
			}),
		ds: ds,
		wc: wc,
	}
	reg.MustRegister(m.siteInfo) // 注册度量指标向量
	return m
}

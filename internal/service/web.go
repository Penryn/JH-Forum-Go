// Copyright 2022 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/servants"
)

var (
	_ Service = (*webService)(nil)
)

// webService 定义了 Web 服务
type webService struct {
	*baseHttpService
}

// Name 返回服务名称
func (s *webService) Name() string {
	return "WebService"
}

// Version 返回服务版本
func (s *webService) Version() *semver.Version {
	return semver.MustParse("v0.5.0")
}

// OnInit 初始化服务，注册服务实例
func (s *webService) OnInit() error {
	s.registerRoute(s, servants.RegisterWebServants)
	return nil
}

// String 返回服务的字符串描述
func (s *webService) String() string {
	return fmt.Sprintf("listen on %s\n", color.GreenString("http://%s:%s", conf.WebServerSetting.HttpIp, conf.WebServerSetting.HttpPort))
}

// newWebEngine 创建一个新的 gin.Engine 实例用于 Web 服务
func newWebEngine() *gin.Engine {
	e := gin.New()
	e.HandleMethodNotAllowed = true
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	// 跨域配置
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	e.Use(cors.New(corsConfig))


	// 默认404
	e.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Not Found",
		})
	})

	// 默认405
	e.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": 405,
			"msg":  "Method Not Allowed",
		})
	})

	return e
}

// newWebService 创建一个新的 WebService 实例
func newWebService() Service {
	addr := conf.WebServerSetting.HttpIp + ":" + conf.WebServerSetting.HttpPort
	server := httpServers.from(addr, func() *httpServer {
		engine := newWebEngine()
		return &httpServer{
			baseServer: newBaseServe(),
			e:          engine,
			server: &http.Server{
				Addr:           addr,
				Handler:        engine,
				ReadTimeout:    conf.WebServerSetting.GetReadTimeout(),
				WriteTimeout:   conf.WebServerSetting.GetWriteTimeout(),
				MaxHeaderBytes: 1 << 20,
			},
		}
	})
	return &webService{
		baseHttpService: &baseHttpService{
			server: server,
		},
	}
}

package debug

import (
	"os"

	"github.com/alimy/tryst/cfg"
	"github.com/pyroscope-io/client/pyroscope"
	"JH-Forum/internal/conf"
	"github.com/sirupsen/logrus"
)

func StartPyroscope() {
	if !cfg.If("Pyroscope") {
		logrus.Infoln("skip Pyroscope because not add Pyroscope feature in config.yaml")
		return
	}
	s := conf.PyroscopeSetting
	c := pyroscope.Config{
		ApplicationName: s.AppName,
		ServerAddress:   s.Endpoint,
		AuthToken:       os.Getenv("PYROSCOPE_AUTH_TOKEN"),
		Logger:          s.GetLogger(),
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
		},
	}
	if len(c.AuthToken) == 0 {
		c.AuthToken = s.AuthToken
	}
	pyroscope.Start(c)
}

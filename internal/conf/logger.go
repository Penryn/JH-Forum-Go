package conf

import (
	"io"


	"github.com/alimy/tryst/cfg"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// newFileLogger 创建一个新的文件日志写入器
func newFileLogger() io.Writer {
	return &lumberjack.Logger{
		Filename:  loggerFileSetting.SavePath + "/" + loggerFileSetting.FileName + loggerFileSetting.FileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}
}

// setupLogger 设置日志系统，根据配置初始化不同的日志输出和钩子
func setupLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{}) // 设置 JSON 格式的日志输出
	logrus.SetLevel(loggerSetting.logLevel())    // 设置日志级别

	// 根据配置注册不同的日志输出和钩子
	cfg.On(cfg.Actions{
		"LoggerFile": func() {
			out := newFileLogger()
			logrus.SetOutput(out)
		},
		"LoggerMeili": func() {
			hook := newMeiliLogHook()
			logrus.SetOutput(io.Discard) // 设置输出到丢弃（不输出到控制台或文件）
			logrus.AddHook(hook)
		},
	})
}

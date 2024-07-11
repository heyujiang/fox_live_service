package config

import (
	"flag"
	"log"
	"os"
	"path"

	"fox_live_service/config/global"
	"fox_live_service/pkg/library/configx"

	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	workPath   = flag.String("workPath", "", "work path")
	configPath = flag.String("configPath", "", "config path")
	configFile = flag.String("configFile", "", "config file path")
)

func init() {
	flag.Parse()
	initEnv(*workPath, *configPath, *configFile)

	initConfig()

	initLog()

	if !global.Config.GetBool("Debug") {
		httpLogPath := path.Join(global.WorkPath, global.LogPath, "http.log") //gin框架的日志
		file, err := os.OpenFile(httpLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		global.AccessLog = file
	}
}

func initEnv(workPath, configPath string, configFile string) {
	if workPath != "" {
		global.WorkPath = workPath
	} else {
		curPath, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		global.WorkPath = curPath
	}
	if configPath != "" {
		global.ConfigPath = configPath
	}
	if configFile != "" {
		global.ConfigFile = configFile
	}
}

func initConfig() {
	global.Config = configx.NewConfigX(global.ConfigPath, global.ConfigCachePrefix, global.ConfigFile)
	global.Config.WatchConfig()
}

func initLog() {
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})

	if !global.Config.GetBool("Debug") {
		logPath := path.Join(global.WorkPath, global.Config.GetString("LogPath"), global.Config.GetString("Name")+".log")

		handler = slog.NewJSONHandler(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    100,
			MaxAge:     20,
			MaxBackups: 30,
			LocalTime:  true,
			Compress:   true,
		}, &slog.HandlerOptions{
			AddSource:   true,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil,
		})
	}
	logger := slog.New(handler)
	slog.NewLogLogger(logger.Handler(), slog.LevelDebug)

	slog.SetDefault(logger)
}

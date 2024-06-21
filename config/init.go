package config

import (
	"flag"
	"log"
	"os"
	"path"

	"fox_live_service/config/global"
	"fox_live_service/pkg/library/configx"

	"github.com/zeromicro/go-zero/core/logx"
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

	logx.MustSetup(logx.LogConf{
		Mode:     global.LogMode,
		Path:     path.Join(global.WorkPath, global.LogPath),
		Encoding: global.LogEncoding,
	})
	logx.Infof("fox-live-server config init...")

	httpLogPath := path.Join(global.WorkPath, global.LogPath, "http.log") //gin框架的日志
	file, err := os.OpenFile(httpLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	global.AccessLog = file
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

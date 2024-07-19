package config

import (
	"flag"
	"fox_live_service/config/global"
	"fox_live_service/pkg/library/configx"
	"log"
	"os"
)

var (
	workPath   = flag.String("workPath", "", "work path")
	configPath = flag.String("configPath", "", "config path")
	configFile = flag.String("configFile", "", "config file path")
)

func init() {
	flag.Parse()
	initEnv(*workPath, *configPath, *configFile)

	global.Config = configx.NewConfigX(global.ConfigPath, global.ConfigFile)
	global.Config.WatchConfig()
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

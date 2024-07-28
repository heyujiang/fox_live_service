package global

import (
	"fox_live_service/pkg/library/configx"
	"os"
)

var (
	WorkPath string

	LogPath = "storages/logs"

	//ConfigPath = "/home/autowise/work/go/fox_live_service/configs"
	ConfigPath = "/Users/fangyamin/go/src/github.com/fox_live_service/configs"
	ConfigFile = "config"
)

var (
	HttpLog *os.File

	Config configx.ConfigI
)

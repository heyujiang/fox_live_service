package global

import (
	"fox_live_service/pkg/library/configx"
	"os"
)

var (
	WorkPath string

	LogPath = "storages/logs"

	ConfigPath        = "/home/autowise/work/go/fox_live_service/config_path"
	ConfigFile        = "config"
	ConfigCachePrefix = "config_"
)

var (
	AccessLog *os.File

	Config configx.ConfigI
)

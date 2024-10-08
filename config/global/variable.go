package global

import (
	"fox_live_service/pkg/library/configx"
	"os"
)

var (
	WorkPath string

	LogPath = "storages/logs"

	//ConfigPath = "/home/autowise/work/go/fox_live_service/configs"
	//C:\Users\syxsx\GolandProjects\fox_live_service\config
	ConfigPath = "C:/Users/syxsx/GolandProjects/fox_live_service/configs"
	ConfigFile = "config"

	FileUploadPath = "storages/upload"
)

var (
	HttpLog *os.File

	Config configx.ConfigI
)

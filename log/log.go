package log

import (
	_ "fox_live_service/config"

	"fox_live_service/config/global"
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
)

func init() {
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

	//处理gin框架日志
	if !global.Config.GetBool("Debug") {
		httpLogPath := path.Join(global.WorkPath, global.LogPath, "http.log") //gin框架的日志
		file, err := os.OpenFile(httpLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		global.HttpLog = file
	}
}

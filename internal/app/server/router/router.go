package router

import (
	"fmt"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/middleware"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Register() *gin.Engine {
	if !global.Config.GetBool("Debug") {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.MultiWriter(global.HttpLog, os.Stdout)
	}

	e := gin.New()

	c := gin.LoggerConfig{
		SkipPaths: []string{""},
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				params.ClientIP,
				params.TimeStamp.Format(time.RFC1123),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
	}
	e.Use(gin.Recovery(), gin.LoggerWithConfig(c), middleware.Cors())

	e.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to jhny manage service")
	})

	registerStaticRouters(e)
	registerUser(e)
	registerProject(e)
	registerNode(e)
	registerUpload(e)
	registerSystem(e)
	registerHome(e)
	registerReport(e)
	return e
}

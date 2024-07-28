package router

import (
	"fmt"
	"fox_live_service/config/global"
	permissions "fox_live_service/internal/app/server/logic/permission"
	"fox_live_service/internal/app/server/middleware"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gin-gonic/gin"
)

func init() {
	policyAdapter, err := gormadapter.NewAdapter(
		"mysql",
		global.Config.GetString("Db.Mysql.DSN"),
		true)
	if err != nil {
		log.Fatalf(err.Error())
	}
	//modelFilePath := "/home/autowise/work/go/fox_live_service/assets/casbin/model.conf"
	modelFilePath := "/Users/fangyamin/go/src/github.com/fox_live_service/assets/casbin/model.conf"
	//modelFilePath := global.Config.GetString("Casbin.ModelFile")

	permission, err := permissions.NewPermissionLogic(func(c *gin.Context) string {
		return c.GetString("username")
	}, modelFilePath, policyAdapter)
	slog.Info("permission init success", permission)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

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

	e.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to fox...")
	})

	registerUser(e)
	registerProject(e)
	registerNode(e)

	return e
}

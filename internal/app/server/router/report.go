package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerReport(e *gin.Engine) {
	reportGroup := e.Group("/api/report", middleware.Auth())
	reportGroup.GET("", handler.ReportHandler.Report)
}

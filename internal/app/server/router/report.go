package router

import (
	"fox_live_service/internal/app/server/handler"
	"github.com/gin-gonic/gin"
)

func registerReport(e *gin.Engine) {
	e.GET("/api/report", handler.ReportHandler.Report)
}

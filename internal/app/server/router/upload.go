package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerUpload(e *gin.Engine) {
	e.POST("/api/upload", middleware.Auth(), handler.UploadHandler.Upload)
}

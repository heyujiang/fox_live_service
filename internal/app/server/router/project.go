package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerProject(e *gin.Engine) {

	project := e.Group("/api/project", middleware.Auth())

	project.GET("/list", handler.ProjectHandler.List)
	project.GET("/:id", handler.ProjectHandler.Info)
	project.POST("", handler.ProjectHandler.Create)
	project.POST("/:id", handler.ProjectHandler.Update)
}

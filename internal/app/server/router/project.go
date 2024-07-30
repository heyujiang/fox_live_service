package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerProject(e *gin.Engine) {

	project := e.Group("/api/project", middleware.Auth())

	project.POST("", handler.ProjectHandler.Create)
	project.DELETE("/:id", handler.ProjectHandler.Delete)
	project.POST("/:id", handler.ProjectHandler.Update)
	project.GET("/:id", handler.ProjectHandler.Info)
	project.GET("", handler.ProjectHandler.List)

	project.GET("/nodes", handler.ProjectNodeHandler.List)

	project.GET("/record", handler.ProjectRecordHandler.List)
	project.POST("/record", handler.ProjectRecordHandler.Create)
	project.DELETE("/record/:id", handler.ProjectRecordHandler.Delete)

	project.POST("/person", handler.ProjectPersonHandler.Create)
	project.DELETE("/person/:id", handler.ProjectPersonHandler.Delete)
	project.GET("/person/:project_id", handler.ProjectPersonHandler.List)

	project.GET("/attached", handler.ProjectAttachedHandler.List)
	project.POST("/attached", handler.ProjectAttachedHandler.Create)
	project.DELETE("/attached/:id", handler.ProjectAttachedHandler.Delete)

}
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
	project.GET("/option", handler.ProjectHandler.Option)

	project.GET("/nodes/:id", handler.ProjectNodeHandler.List)
	project.GET("/nodes/option/:id", handler.ProjectNodeHandler.Option)

	project.GET("/record", handler.ProjectRecordHandler.List)
	project.POST("/record", handler.ProjectRecordHandler.Create)
	project.POST("/record/:id", handler.ProjectRecordHandler.Update)
	project.DELETE("/record/:id", handler.ProjectRecordHandler.Delete)

	project.POST("/person", handler.ProjectPersonHandler.Create)
	project.DELETE("/person/:id", handler.ProjectPersonHandler.Delete)
	project.GET("/person/:projectId", handler.ProjectPersonHandler.List)

	project.POST("/contact", handler.ProjectContactHandler.Create)
	project.DELETE("/contact/:id", handler.ProjectContactHandler.Delete)
	project.GET("/contact/:projectId", handler.ProjectContactHandler.List)

	project.GET("/attached", handler.ProjectAttachedHandler.List)
	project.POST("/attached", handler.ProjectAttachedHandler.Create)
	project.DELETE("/attached/:id", handler.ProjectAttachedHandler.Delete)

}

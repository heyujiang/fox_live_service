package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerNode(e *gin.Engine) {
	node := e.Group("/api/node", middleware.Auth())

	node.GET("/list", handler.ProjectHandler.List)
	node.GET("/:id", handler.ProjectHandler.Info)
	node.POST("", handler.ProjectHandler.Create)
	node.POST("/:id", handler.ProjectHandler.Update)
}

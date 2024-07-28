package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerNode(e *gin.Engine) {

	node := e.Group("/api/node", middleware.Auth())

	node.POST("", handler.NodeHandler.Create)
	node.DELETE("/:id", handler.NodeHandler.Delete)
	node.POST("/:id", handler.NodeHandler.Update)
	node.GET("/:id", handler.NodeHandler.Info)
	node.GET("", handler.NodeHandler.List)
	node.GET("/parent", handler.NodeHandler.Parent)
}

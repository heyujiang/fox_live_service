package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerHome(e *gin.Engine) {
	home := e.Group("/api/home").Use(middleware.Auth())

	home.GET("/viewCount", handler.HomeHandler.ViewCount)
	home.GET("/latestRecord", handler.HomeHandler.LatestRecord)
	home.GET("/latestProject", handler.HomeHandler.LatestProject)
	home.GET("/personCapacity", handler.HomeHandler.PersonCapacity)
	home.GET("/userData", handler.HomeHandler.UserData)
}

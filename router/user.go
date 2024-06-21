package router

import (
	"fox_live_service/handler"
	"github.com/gin-gonic/gin"
)

func registerUser(e *gin.Engine) {
	e.GET("/login", handler.UserHandler.Login)
}

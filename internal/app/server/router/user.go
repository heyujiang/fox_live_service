package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerUser(e *gin.Engine) {
	e.POST("/login", handler.UserHandler.Login) // 登录

	user := e.Group("/api/user").Use(middleware.Auth())
	user.GET("/logout/:id", handler.UserHandler.Logout) // 退出

	user.GET("/list", handler.UserHandler.List)             // 用户信息
	user.GET("/:id", handler.UserHandler.Info)              // 用户信息
	user.POST("/", handler.UserHandler.Create)              // 添加用户
	user.POST("/:id", handler.UserHandler.Update)           // 修改用户
	user.POST("/enable/:id", handler.UserHandler.Enabled)   // 启用用户
	user.POST("/disable/:id", handler.UserHandler.Disabled) // 禁用用户
}

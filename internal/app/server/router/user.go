package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerUser(e *gin.Engine) {
	e.POST("/login", handler.UserHandler.Login) // 登录

	user := e.Group("/api/user").Use(middleware.Auth())
	user.POST("/logout", handler.UserHandler.Logout) // 退出

	user.GET("", handler.UserHandler.List)                 // 用户信息
	user.GET("/:id", handler.UserHandler.Info)             // 用户信息
	user.POST("", handler.UserHandler.Create)              // 添加用户
	user.POST("/:id", handler.UserHandler.Update)          // 修改用户
	user.POST("/enable/:id", handler.UserHandler.Enable)   // 启用用户
	user.POST("/disable/:id", handler.UserHandler.Disable) // 禁用用户
	user.DELETE("/:id", handler.UserHandler.Delete)        // 删除用户
	user.GET("/options", handler.UserHandler.Options)      // 删除用户

	user.GET("/menus", handler.UserHandler.Menus) //用户菜单

	user.GET("/getUserInfo", handler.UserHandler.GetUserInfo)        // 获取登录用户信息
	user.POST("/updateAvatar", handler.UserHandler.UpdateAvatar)     // 修改头像
	user.POST("/updateBasic", handler.UserHandler.UpdateBasic)       // 修改基础信息
	user.POST("/updatePassword", handler.UserHandler.UpdatePassword) // 修改密码
}

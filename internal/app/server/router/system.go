package router

import (
	"fox_live_service/internal/app/server/handler"
	"fox_live_service/internal/app/server/middleware"
	"github.com/gin-gonic/gin"
)

func registerSystem(e *gin.Engine) {

	system := e.Group("/api/system", middleware.Auth())

	system.POST("/rule", handler.RuleHandler.Create)
	system.DELETE("/rule/:id", handler.RuleHandler.Delete)
	system.POST("/rule/:id", handler.RuleHandler.Update)
	system.POST("/rule/updateStatus/:id", handler.RuleHandler.UpdateStatus)
	system.GET("/rule", handler.RuleHandler.List)
	system.GET("/rule/parents", handler.RuleHandler.Parents)
	system.GET("/rule/getRules/:id", handler.RuleHandler.GetRoleRules)

	system.POST("/role", handler.RoleHandler.Create)
	system.DELETE("/role/:id", handler.RoleHandler.Delete)
	system.POST("/role/:id", handler.RoleHandler.Update)
	system.POST("/role/updateStatus/:id", handler.RoleHandler.UpdateStatus)
	system.GET("/role", handler.RoleHandler.List)
	system.GET("/role/parents", handler.RoleHandler.Parents)

	system.POST("/dept", handler.DeptHandler.Create)
	system.DELETE("/dept/:id", handler.DeptHandler.Delete)
	system.POST("/dept/:id", handler.DeptHandler.Update)
	system.POST("/dept/updateStatus/:id", handler.DeptHandler.UpdateStatus)
	system.GET("/dept", handler.DeptHandler.List)
	system.GET("/dept/parents", handler.DeptHandler.Parents)

}

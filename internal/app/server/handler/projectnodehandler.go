package handler

import "github.com/gin-gonic/gin"

var ProjectNodeHandler = newProjectNodeHandler()

type projectNodeHandler struct{}

func newProjectNodeHandler() *projectNodeHandler {
	return &projectNodeHandler{}
}

// List 项目节点记录列表
func (h *projectNodeHandler) List(c *gin.Context) {}

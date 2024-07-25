package handler

import "github.com/gin-gonic/gin"

var ProjectAttachedHandler = newProjectAttachedHandler()

type projectAttachedHandler struct{}

func newProjectAttachedHandler() *projectAttachedHandler {
	return &projectAttachedHandler{}
}

// List 项目节点记录附件列表
func (p *projectAttachedHandler) List(c *gin.Context) {}

// Create 创建项目节点记录附件
func (p *projectAttachedHandler) Create(c *gin.Context) {}

// Delete 删除项目节点记录附件
func (p *projectAttachedHandler) Delete(c *gin.Context) {}

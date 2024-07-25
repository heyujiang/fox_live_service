package handler

import "github.com/gin-gonic/gin"

var ProjectRecordHandler = newProjectRecordHandler()

type projectRecordHandler struct{}

func newProjectRecordHandler() *projectRecordHandler {
	return &projectRecordHandler{}
}

// List 项目节点记录列表
func (p *projectRecordHandler) List(c *gin.Context) {}

// Create 创建项目节点记录
func (p *projectRecordHandler) Create(c *gin.Context) {}

// Delete 删除项目节点记录
func (p *projectRecordHandler) Delete(c *gin.Context) {}

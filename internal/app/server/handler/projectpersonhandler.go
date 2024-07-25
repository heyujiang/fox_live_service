package handler

import "github.com/gin-gonic/gin"

var ProjectPersonHandler = newProjectPersonHandler()

type projectPersonHandler struct{}

func newProjectPersonHandler() *projectPersonHandler {
	return &projectPersonHandler{}
}

// List 项目成员列表
func (p *projectPersonHandler) List(c *gin.Context) {}

// Create 创建项目成员
func (p *projectPersonHandler) Create(c *gin.Context) {}

// Delete 删除项目成员
func (p *projectPersonHandler) Delete(c *gin.Context) {}

package handler

import "github.com/gin-gonic/gin"

var ProjectHandler = newProjectHandler()

type projectHandler struct{}

func newProjectHandler() *projectHandler {
	return &projectHandler{}
}

// List 项目列表
func (p *projectHandler) List(c *gin.Context) {

}

// Info 项目详情
func (p *projectHandler) Info(c *gin.Context) {

}

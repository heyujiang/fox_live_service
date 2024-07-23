package handler

import "github.com/gin-gonic/gin"

var ProjectHandler = newProjectHandler()

type projectHandler struct{}

func newProjectHandler() *projectHandler {
	return &projectHandler{}
}

func (p *projectHandler) List(c *gin.Context) {

}

func (p *projectHandler) Info(c *gin.Context) {

}

func (p *projectHandler) Create(c *gin.Context) {

}

func (p *projectHandler) Update(c *gin.Context) {

}

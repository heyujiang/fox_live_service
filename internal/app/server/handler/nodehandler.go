package handler

import "github.com/gin-gonic/gin"

var NodeHandler = newNodeHandler()

type nodeHandler struct{}

func newNodeHandler() *nodeHandler {
	return &nodeHandler{}
}

func (n *nodeHandler) List(c *gin.Context) {}

func (n *nodeHandler) Info(c *gin.Context) {}

func (n *nodeHandler) Create(c *gin.Context) {}

func (n *nodeHandler) Update(c *gin.Context) {}

func (n *nodeHandler) Delete(c *gin.Context) {}

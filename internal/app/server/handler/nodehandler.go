package handler

import (
	"fox_live_service/internal/app/server/logic/node"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var NodeHandler = newNodeHandler()

type nodeHandler struct{}

func newNodeHandler() *nodeHandler {
	return &nodeHandler{}
}

func (n *nodeHandler) List(c *gin.Context) {
	res, err := node.BisLogic.List()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (n *nodeHandler) Info(c *gin.Context) {
	var req node.ReqNodeInfo
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := node.BisLogic.Info(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (n *nodeHandler) Create(c *gin.Context) {
	var req node.ReqCreateNode
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := node.BisLogic.Create(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (n *nodeHandler) Update(c *gin.Context) {
	var reqUri node.ReqUriUpdateNode
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody node.ReqBodyUpdateNode
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := node.ReqUpdateNode{
		ReqUriUpdateNode:  reqUri,
		ReqBodyUpdateNode: reqBody,
	}

	res, err := node.BisLogic.Update(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (n *nodeHandler) Delete(c *gin.Context) {
	var req node.ReqDeleteNode
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := node.BisLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (n *nodeHandler) Parent(c *gin.Context) {
	res, err := node.BisLogic.Parent()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (n *nodeHandler) Options(c *gin.Context) {
	res, err := node.BisLogic.Options()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

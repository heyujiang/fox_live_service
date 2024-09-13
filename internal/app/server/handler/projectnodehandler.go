package handler

import (
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var ProjectNodeHandler = newProjectNodeHandler()

type projectNodeHandler struct{}

func newProjectNodeHandler() *projectNodeHandler {
	return &projectNodeHandler{}
}

func (h *projectNodeHandler) Create(c *gin.Context) {
	var req project.ReqCreateProjectNode
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.NodeLogic.Create(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectNodeHandler) Delete(c *gin.Context) {
	var req project.ReqDeleteProjectNode
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.NodeLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectNodeHandler) Update(c *gin.Context) {
	var reqUri project.ReqUriUpdateProjectNode
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody project.ReqBodyUpdateProjectNode
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := project.ReqUpdateProjectNode{
		ReqUriUpdateProjectNode:  reqUri,
		ReqBodyUpdateProjectNode: reqBody,
	}

	res, err := project.NodeLogic.Update(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectNodeHandler) Info(c *gin.Context) {
	var req project.ReqInfoProjectNode
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.NodeLogic.Info(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectNodeHandler) List(c *gin.Context) {
	var req project.ReqProjectNodeList
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.NodeLogic.List(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectNodeHandler) Option(c *gin.Context) {
	var req project.ReqProjectNodeOption
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.NodeLogic.Option(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

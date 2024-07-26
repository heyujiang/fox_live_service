package handler

import (
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var ProjectAttachedHandler = newProjectAttachedHandler()

type projectAttachedHandler struct{}

func newProjectAttachedHandler() *projectAttachedHandler {
	return &projectAttachedHandler{}
}

func (h *projectAttachedHandler) Create(c *gin.Context) {
	var req project.ReqCreateProjectAttached
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.AttachedLogic.Create(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectAttachedHandler) Delete(c *gin.Context) {
	var req project.ReqDeleteProjectAttached
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.AttachedLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectAttachedHandler) Update(c *gin.Context) {
	var reqUri project.ReqUriUpdateProjectAttached
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody project.ReqBodyUpdateProjectAttached
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := project.ReqUpdateProjectAttached{
		ReqUriUpdateProjectAttached:  reqUri,
		ReqBodyUpdateProjectAttached: reqBody,
	}

	res, err := project.AttachedLogic.Update(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectAttachedHandler) Info(c *gin.Context) {
	var req project.ReqInfoProjectAttached
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.AttachedLogic.Info(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectAttachedHandler) List(c *gin.Context) {
	var req project.ReqProjectAttachedList
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.AttachedLogic.List(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

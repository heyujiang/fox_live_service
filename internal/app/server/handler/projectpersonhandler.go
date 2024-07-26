package handler

import (
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var ProjectPersonHandler = newProjectPersonHandler()

type projectPersonHandler struct{}

func newProjectPersonHandler() *projectPersonHandler {
	return &projectPersonHandler{}
}

func (h *projectPersonHandler) Create(c *gin.Context) {
	var req project.ReqCreateProjectPerson
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.PersonLogic.Create(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectPersonHandler) Delete(c *gin.Context) {
	var req project.ReqDeleteProjectPerson
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.PersonLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectPersonHandler) Update(c *gin.Context) {
	var reqUri project.ReqUriUpdateProjectPerson
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody project.ReqBodyUpdateProjectPerson
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := project.ReqUpdateProjectPerson{
		ReqUriUpdateProjectPerson:  reqUri,
		ReqBodyUpdateProjectPerson: reqBody,
	}

	res, err := project.PersonLogic.Update(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectPersonHandler) Info(c *gin.Context) {
	var req project.ReqInfoProjectPerson
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.PersonLogic.Info(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectPersonHandler) List(c *gin.Context) {
	var req project.ReqProjectPersonList
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.PersonLogic.List(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

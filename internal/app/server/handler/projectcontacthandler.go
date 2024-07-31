package handler

import (
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var ProjectContactHandler = newProjectContactHandler()

type projectContactHandler struct{}

func newProjectContactHandler() *projectContactHandler {
	return &projectContactHandler{}
}

func (h *projectContactHandler) Create(c *gin.Context) {
	var req project.ReqCreateProjectContact
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.ContactLogic.Create(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectContactHandler) Delete(c *gin.Context) {
	var req project.ReqDeleteProjectContact
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.ContactLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectContactHandler) List(c *gin.Context) {
	var req project.ReqProjectContactList
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.ContactLogic.List(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

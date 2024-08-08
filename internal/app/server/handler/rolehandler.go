package handler

import (
	"fox_live_service/internal/app/server/logic/system"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var RoleHandler = newRoleHandler()

type roleHandler struct{}

func newRoleHandler() *roleHandler {
	return &roleHandler{}
}

func (h *roleHandler) Create(c *gin.Context) {
	var req system.ReqCreateRole
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := system.RoleLogic.Create(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *roleHandler) Delete(c *gin.Context) {
	var req system.ReqDeleteRole
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := system.RoleLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *roleHandler) Update(c *gin.Context) {
	var reqUri system.ReqUpdateRoleUri
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody system.ReqUpdateRoleBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := system.ReqUpdateRole{
		ReqUpdateRoleUri:  reqUri,
		ReqUpdateRoleBody: reqBody,
	}

	res, err := system.RoleLogic.Update(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *roleHandler) List(c *gin.Context) {
	res, err := system.RoleLogic.List()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *roleHandler) UpdateStatus(c *gin.Context) {
	var reqUri system.ReqUpdateRoleUri
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody system.ReqUpdateRoleStatusBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := system.ReqUpdateRoleStatus{
		ReqUpdateRoleUri:        reqUri,
		ReqUpdateRoleStatusBody: reqBody,
	}

	res, err := system.RoleLogic.UpdateStatus(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *roleHandler) Parents(c *gin.Context) {
	res, err := system.RoleLogic.Parents()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

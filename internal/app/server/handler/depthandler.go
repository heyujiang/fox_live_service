package handler

import (
	"fox_live_service/internal/app/server/logic/system"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var DeptHandler = newDeptHandler()

type deptHandler struct{}

func newDeptHandler() *deptHandler {
	return &deptHandler{}
}

func (h *deptHandler) Create(c *gin.Context) {
	var req system.ReqCreateDept
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := system.DeptLogic.Create(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *deptHandler) Delete(c *gin.Context) {
	var req system.ReqDeleteDept
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := system.DeptLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *deptHandler) Update(c *gin.Context) {
	var reqUri system.ReqUpdateDeptUri
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody system.ReqUpdateDeptBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := system.ReqUpdateDept{
		ReqUpdateDeptUri:  reqUri,
		ReqUpdateDeptBody: reqBody,
	}

	res, err := system.DeptLogic.Update(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *deptHandler) List(c *gin.Context) {
	res, err := system.DeptLogic.List()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *deptHandler) UpdateStatus(c *gin.Context) {
	var reqUri system.ReqUpdateDeptUri
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody system.ReqUpdateDeptStatusBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := system.ReqUpdateDeptStatus{
		ReqUpdateDeptUri:        reqUri,
		ReqUpdateDeptStatusBody: reqBody,
	}

	res, err := system.DeptLogic.UpdateStatus(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *deptHandler) Parents(c *gin.Context) {
	res, err := system.DeptLogic.Parents()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

package handler

import (
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var ProjectRecordHandler = newProjectRecordHandler()

type projectRecordHandler struct{}

func newProjectRecordHandler() *projectRecordHandler {
	return &projectRecordHandler{}
}

func (h *projectRecordHandler) Create(c *gin.Context) {
	var req project.ReqCreateProjectRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.RecordLogic.Create(&req, c.GetInt("uid"), c.GetString("username"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectRecordHandler) Delete(c *gin.Context) {
	var req project.ReqDeleteProjectRecord
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.RecordLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectRecordHandler) Update(c *gin.Context) {
	var reqUri project.ReqUriUpdateProjectRecord
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody project.ReqBodyUpdateProjectRecord
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := project.ReqUpdateProjectRecord{
		ReqUriUpdateProjectRecord:  reqUri,
		ReqBodyUpdateProjectRecord: reqBody,
	}

	res, err := project.RecordLogic.Update(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectRecordHandler) Info(c *gin.Context) {
	var req project.ReqInfoProjectRecord
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.RecordLogic.Info(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectRecordHandler) List(c *gin.Context) {
	var req project.ReqProjectRecordList
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	if !c.GetBool("isSuper") {
		req.UserId = c.GetInt("uid")
	}
	res, err := project.RecordLogic.List(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectRecordHandler) ListNoPage(c *gin.Context) {
	var req project.ReqProjectRecordList
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.RecordLogic.ListNoPage(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectRecordHandler) GetLatestRecords(c *gin.Context) {
	res, err := project.RecordLogic.GetLatestRecords(c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectRecordHandler) GetTeams(c *gin.Context) {
	res, err := project.RecordLogic.GetTeams(c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

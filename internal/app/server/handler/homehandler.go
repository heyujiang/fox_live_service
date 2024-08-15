package handler

import (
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var HomeHandler = newHomeHandler()

type homeHandler struct{}

func newHomeHandler() *homeHandler {
	return &homeHandler{}
}

// ViewCount 首页数据汇总
func (h *homeHandler) ViewCount(c *gin.Context) {
	res, err := project.BisLogic.ViewCount()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

// LatestRecord 项目最新进展
func (h *homeHandler) LatestRecord(c *gin.Context) {
	res, err := project.RecordLogic.GetAllLatestRecords()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}
	common.ResponseOK(c, res)
	return
}

func (h *homeHandler) LatestProject(c *gin.Context) {
	var req project.ReqGetLatestProject
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	res, err := project.BisLogic.GetLatestProject(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}
	common.ResponseOK(c, res)
	return
}

func (h *homeHandler) PersonCapacity(c *gin.Context) {
	res, err := project.BisLogic.PersonCapacity()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}
	common.ResponseOK(c, res)
	return
}

package handler

import (
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/pkg/common"
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

package handler

import (
	"errors"
	"fox_live_service/internal/app/server/logic/report"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var ReportHandler = newReportHandler()

type reportHandler struct{}

func newReportHandler() *reportHandler {
	return &reportHandler{}
}

func (r *reportHandler) Report(c *gin.Context) {
	var req report.ReqReport
	if err := c.ShouldBindQuery(&req); err != nil {
		var errs validator.ValidationErrors
		if !errors.As(err, &errs) {
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
			return
		} else {
			common.ResponseErr(c, errorx.ValidationTran(errs))
			return
		}
	}

	if len(req.TimeRange) != 2 {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "请选择时间范围"))
		return
	}

	if !c.GetBool("isSuper") && !c.GetBool("isSystem") { //非超级管理员和系统账户只能看到自己的数据
		req.UserId = c.GetInt("uid")
	} else {
		if req.UserId == 0 {
			req.UserId = c.GetInt("uid")
		}
	}

	res, err := report.ReportLogic.Report(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

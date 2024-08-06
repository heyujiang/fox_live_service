package handler

import (
	"fox_live_service/internal/app/server/logic/upload"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var (
	UploadHandler = uploadHandler{}
)

type uploadHandler struct {
}

func newUploadHandler() *uploadHandler {
	return &uploadHandler{}
}

// Upload 文件上传接口
func (h *uploadHandler) Upload(c *gin.Context) {
	var req upload.ReqFileUpload
	if err := c.ShouldBind(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	resp, err := upload.BisLogic.Upload(c, &req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, resp)
	return
}

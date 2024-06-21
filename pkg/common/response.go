package common

import (
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func ResponseErr(c *gin.Context, err error) {
	errX := errorx.ParseErrorX(err)
	c.JSON(http.StatusOK, Response{
		Code: errX.GetCode(),
		Msg:  errX.GetMsg(),
		Data: nil,
	})
}

func ResponseErrData(c *gin.Context, data interface{}, err error) {
	errX := errorx.ParseErrorX(err)
	c.JSON(http.StatusOK, Response{
		Code: errX.GetCode(),
		Msg:  errX.GetMsg(),
		Data: data,
	})
}

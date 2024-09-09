package handler

import (
	"bytes"
	"fmt"
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var ProjectHandler = newProjectHandler()

type projectHandler struct{}

func newProjectHandler() *projectHandler {
	return &projectHandler{}
}

func (h *projectHandler) Create(c *gin.Context) {
	var req project.ReqCreateProject
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.BisLogic.Create(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectHandler) Delete(c *gin.Context) {
	var req project.ReqDeleteProject
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.BisLogic.Delete(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectHandler) Update(c *gin.Context) {
	var reqUri project.ReqUriUpdateProject
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody project.ReqBodyUpdateProject
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := project.ReqUpdateProject{
		ReqUriUpdateProject:  reqUri,
		ReqBodyUpdateProject: reqBody,
	}

	res, err := project.BisLogic.Update(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectHandler) Info(c *gin.Context) {
	var req project.ReqInfoProject
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.BisLogic.Info(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectHandler) List(c *gin.Context) {
	var req project.ReqProjectList
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.BisLogic.List(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectHandler) Option(c *gin.Context) {
	var req project.ReqProjectOption
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.BisLogic.Option(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectHandler) GetMyProject(c *gin.Context) {
	res, err := project.BisLogic.GetMyProject(c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *projectHandler) Export(c *gin.Context) {
	var req project.ReqFromProjectList
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := project.BisLogic.Export(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	c.Header("Content-Length", strconv.Itoa(len(res.Data.Bytes())))
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename*=UTF-8''%s", url.QueryEscape("project.xlsx")))
	c.Header("Content-Type", "application/octet-stream")

	//xlsx，设置后缀为xlsx类型的表格
	//c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	r := bytes.NewReader(res.Data.Bytes())
	// 返回给浏览器
	http.ServeContent(c.Writer, c.Request, "project.xlsx", time.Now(), r)
	return
}

package handler

import (
	"fox_live_service/internal/app/server/logic/user"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var UserHandler = newUserHandler()

type userHandler struct{}

func newUserHandler() *userHandler {
	return &userHandler{}
}

// Login 用户名密码登录
func (u *userHandler) Login(c *gin.Context) {
	var req user.ReqLogin
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := user.LoginLogic.Login(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return

}

func (u *userHandler) Logout(c *gin.Context) {

}

// List 用户列表
func (u *userHandler) List(c *gin.Context) {
	var req user.ReqUserList
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := user.BisLogic.List(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) Info(c *gin.Context) {

}

// Create 创建用户
func (u *userHandler) Create(c *gin.Context) {
	var req user.ReqCreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := user.BisLogic.Create(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) Update(c *gin.Context) {

}

func (u *userHandler) Enabled(c *gin.Context) {

}

func (u *userHandler) Disabled(c *gin.Context) {

}

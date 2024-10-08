package handler

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/logic/user"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var UserHandler = newUserHandler()

type userHandler struct{}

func newUserHandler() *userHandler {
	return &userHandler{}
}

// Login 用户名密码登录
func (u *userHandler) Login(c *gin.Context) {
	var req user.ReqLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := user.LoginLogic.Login(&req, global.Config.GetString("AccessToken.JwtTokenKey"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

// Logout 用户退出
func (u *userHandler) Logout(c *gin.Context) {
	common.ResponseOK(c, nil)
	return
}

// List 用户列表
func (u *userHandler) List(c *gin.Context) {
	var req user.ReqUserList
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

	res, err := user.BisLogic.List(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) Info(c *gin.Context) {
	var req user.ReqUserInfo
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := user.BisLogic.Info(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

// Create 创建用户
func (u *userHandler) Create(c *gin.Context) {
	var req user.ReqCreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if !errors.As(err, &errs) {
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
			return
		} else {
			common.ResponseErr(c, errorx.ValidationTran(errs))
			return
		}
	}

	res, err := user.BisLogic.Create(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

// Update 修改用户
func (u *userHandler) Update(c *gin.Context) {
	var reqUri user.ReqUriUpdateUser
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	var reqBody user.ReqBodyUpdateUser
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := user.ReqUpdateUser{
		ReqUriUpdateUser:  reqUri,
		ReqBodyUpdateUser: reqBody,
	}

	res, err := user.BisLogic.Update(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) Delete(c *gin.Context) {
	var req user.ReqDeleteUser
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := user.BisLogic.Delete(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) Enable(c *gin.Context) {
	var req user.ReqEnableUser
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := user.BisLogic.Enable(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) Disable(c *gin.Context) {
	var req user.ReqDisableUser
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := user.BisLogic.Disable(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) Menus(c *gin.Context) {
	res, _ := user.MenuLogic.GetMenus(c.GetInt("uid"))
	common.ResponseOK(c, res)
	return
}

func (u *userHandler) GetUserInfo(c *gin.Context) {
	res, err := user.AccountLogic.UserInfo(c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}
	common.ResponseOK(c, res)
	return
}

func (u *userHandler) Options(c *gin.Context) {
	res, _ := user.BisLogic.Options()
	common.ResponseOK(c, res)
	return
}

func (u *userHandler) UpdateBasic(c *gin.Context) {
	var req user.ReqUpdateBasic
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	res, err := user.AccountLogic.UpdateBasic(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) UpdateAvatar(c *gin.Context) {
	var req user.ReqUpdateAvatar
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	res, err := user.AccountLogic.UpdateAvatar(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (u *userHandler) UpdatePassword(c *gin.Context) {
	var req user.ReqUpdatePassword
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	res, err := user.AccountLogic.UpdatePassword(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

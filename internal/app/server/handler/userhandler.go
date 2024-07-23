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

func (u *userHandler) Login(c *gin.Context) {
	user.BisLogic.Login()
}

func (u *userHandler) Logout(c *gin.Context) {

}

// List 用户列表
func (u *userHandler) List(c *gin.Context) {
	var req user.ReqUserList
	if err := c.ShouldBindQuery(&req); err != nil {
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

func (u *userHandler) Create(c *gin.Context) {

}

func (u *userHandler) Update(c *gin.Context) {

}

func (u *userHandler) Enabled(c *gin.Context) {

}

func (u *userHandler) Disabled(c *gin.Context) {

}

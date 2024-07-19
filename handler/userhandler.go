package handler

import (
	"fox_live_service/logic"
	"github.com/gin-gonic/gin"
)

var UserHandler = newUserHandler()

type userHandler struct{}

func newUserHandler() *userHandler {
	return &userHandler{}
}

func (u *userHandler) Login(ctx *gin.Context) {
	logic.UserLogic.Login()
}

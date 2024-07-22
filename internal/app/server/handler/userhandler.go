package handler

import (
	"fox_live_service/internal/app/server/logic"
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

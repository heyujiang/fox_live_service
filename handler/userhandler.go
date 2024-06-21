package handler

import (
	"fox_live_service/logic"
	"github.com/gin-gonic/gin"
	"time"
)

var UserHandler = newUserHandler()

func newUserHandler() *userHandler {
	return &userHandler{
		now: time.Now(),
		lg:  logic.ULogic,
	}
}

type userHandler struct {
	now time.Time
	lg  *logic.UserLogic
}

func (u *userHandler) Login(ctx *gin.Context) {
	//fmt.Println("handler : ", u.now)

	//u.lg.Login()

	gg := logic.NewUserLogic()
	gg.Login()
}

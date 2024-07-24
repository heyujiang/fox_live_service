package user

import (
	"errors"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
)

var LoginLogic = newLogicLogic()

type (
	loginLogic struct{}

	ReqLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	RespLogin struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Token    string `json:"token"`
	}
)

func newLogicLogic() *loginLogic {
	return &loginLogic{}
}

// Login 用户名密码登录
func (l loginLogic) Login(req *ReqLogin) (*RespLogin, error) {
	//查询用户是否存在
	user, err := model.UserModel.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrUserNotExist, "用户名不存在或密码错误")
		}
		slog.Error("login error")
		return nil, err
	}
	//判断密码是否争正确
	if user.Username != req.Username || user.Password != req.Password {
		slog.Error("password error")
		return nil, errorx.NewErrorX(errorx.ErrUserNotExist, "用户名不存在或密码错误")
	}

	//生成TOKEN
	token := "thisisusernamelogintoken"
	return &RespLogin{
		Id:       user.Id,
		Username: user.Username,
		Token:    token,
	}, nil
}

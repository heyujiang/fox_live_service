package user

import (
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
)

var AccountLogic = newAccountLogic()

type (
	accountLogic struct{}

	RespAccountUserInfo struct {
		UserId       int    `json:"userId"`
		Username     string `json:"username"`
		Name         string `json:"name"`
		Nickname     string `json:"nickname"`
		Avatar       string `json:"avatar"`
		Role         string `json:"role"`
		Introduction string `json:"introduction"`
		Email        string `json:"email"`
	}
)

func newAccountLogic() *accountLogic {
	return &accountLogic{}
}

func (a *accountLogic) UserInfo(uid int) (*RespAccountUserInfo, error) {
	//查询用户名是否存在
	user, err := model.UserModel.Find(uid)
	if err != nil {
		slog.Error("get account user info error ： ", "uid", uid, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取账户信息失败")
	}

	return &RespAccountUserInfo{
		UserId:       user.Id,
		Username:     user.Username,
		Name:         user.Name,
		Nickname:     user.NickName,
		Avatar:       user.Avatar,
		Email:        user.Email,
		Role:         "",
		Introduction: "",
	}, nil

}

package logic

import (
	"fox_live_service/model"
	"golang.org/x/exp/slog"
)

var UserLogic = newUserLogic()

func newUserLogic() *userLogic {
	return &userLogic{}
}

type userLogic struct{}

func (ul *userLogic) Login() {
	if err := model.UserModel.Insert(&model.User{
		Username:    "jiangyu",
		Password:    "1324321",
		PhoneNumber: "15658086185",
		Email:       "1393870072@qq.com",
		Name:        "江屿",
		NickName:    "He Y J",
		Avatar:      "avatar",
		Model: model.Model{
			CreateId: 1,
			UpdateId: 1,
		},
	}); err != nil {
		slog.Error(err.Error())
	}
	return
}

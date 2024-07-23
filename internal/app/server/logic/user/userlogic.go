package user

import (
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/logic"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
)

var BisLogic = newBisLogic()

func newBisLogic() *bisLogic {
	return &bisLogic{}
}

type (
	bisLogic struct{}

	ReqUserList struct {
		logic.ReqPage
		ReqBodyUserList
	}

	ReqBodyUserList struct {
		Id          uint   `form:"id"`
		PhoneNumber string `form:"phone_number"`
		Name        string `form:"name"`
		State       int    `form:"state"`
	}

	RespUserList struct {
		Page  int     `json:"page"`
		Size  int     `json:"size"`
		Count int     `json:"count"`
		List  []*Item `json:"list"`
	}

	Item struct {
		Id          uint   `json:"id"`
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		Avatar      string `json:"avatar"`
		State       int    `json:"state"`
		CreatedId   uint   `json:"create_id"`
		UpdatedId   uint   `json:"update_id"`
		CreatedAt   string `json:"create_at"`
		UpdatedAT   string `json:"update_at"`
	}
)

func (b *bisLogic) Login() {
	if err := model.UserModel.Insert(&model.User{
		Username:    "jiangyu",
		Password:    "1324321",
		PhoneNumber: "15658086185",
		Email:       "1393870072@qq.com",
		Name:        "江屿",
		NickName:    "He Y J",
		Avatar:      "avatar",
		State:       model.UserStatusEnable,
		CreatedId:   1,
		UpdatedId:   1,
	}); err != nil {
		slog.Error(err.Error())
	}
	return
}

func (b *bisLogic) List(req *ReqUserList) (*RespUserList, error) {
	cond := b.buildSearchCond(req)
	totalCount, err := model.UserModel.GetUserCountCond(cond)
	if err != nil {
		slog.Error("list user get user count error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取用户列表错误")
	}
	users, err := model.UserModel.GetUsersByCond(cond, req.Page, req.Size)
	if err != nil {
		slog.Error("list user get user list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取用户列表错误")
	}

	items := make([]*Item, 0, len(users))
	for _, user := range users {
		items = append(items, &Item{
			Id:          user.Id,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Name:        user.Name,
			Avatar:      user.Avatar,
			State:       user.State,
			CreatedId:   user.CreatedId,
			UpdatedId:   user.UpdatedId,
			CreatedAt:   user.CreatedAt.Format(global.TimeFormat),
			UpdatedAT:   user.CreatedAt.Format(global.TimeFormat),
		})
	}

	return &RespUserList{
		Page:  req.Page,
		Size:  req.Size,
		Count: totalCount,
		List:  items,
	}, nil
}

func (b *bisLogic) buildSearchCond(req *ReqUserList) *model.UserCond {
	return nil
}

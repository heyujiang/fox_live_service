package user

import (
	"crypto/md5"
	"errors"
	"fmt"
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
		Id          int    `form:"id"`
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
		Id          int    `json:"id"`
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		Avatar      string `json:"avatar"`
		State       int    `json:"state"`
		CreatedId   int    `json:"create_id"`
		UpdatedId   int    `json:"update_id"`
		CreatedAt   string `json:"create_at"`
		UpdatedAT   string `json:"update_at"`
	}

	ReqCreateUser struct {
		Username    string `json:"username"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}

	RespCreateUser struct {
		Item
	}
)

// Create 创建用户
func (b *bisLogic) Create(req *ReqCreateUser, uid int) (*RespCreateUser, error) {
	//查询用户名是否存在
	user, err := model.UserModel.FindByUsername(req.Username)
	if err != nil && !errors.Is(err, model.ErrNotRecord) {
		slog.Error("create user error ： get user by username error ", "username", req.Username, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建用户失败")
	}
	if user != nil {
		slog.Error("create user error ： username is exist", "username", req.Username)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "用户名已经存在")
	}

	//查询手机号是否存在
	user, err = model.UserModel.FindByPhoneNumber(req.PhoneNumber)
	if err != nil && !errors.Is(err, model.ErrNotRecord) {
		slog.Error("create user error ： get user by phone number error ", "phone number", req.PhoneNumber, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建用户失败")
	}
	if user != nil {
		slog.Error("create user error ：  phone number  is exist", " phone number ", req.PhoneNumber)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "手机号已经存在")
	}

	insertUser := model.User{
		Username:    req.Username,
		Password:    md5Password(req.Password),
		PhoneNumber: req.PhoneNumber,
		Email:       "",
		Name:        req.Name,
		Avatar:      "",
		State:       model.UserStatusEnable,
		CreatedId:   uid,
		UpdatedId:   uid,
	}

	if err := model.UserModel.Insert(&insertUser); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建用户失败")
	}

	return &RespCreateUser{
		Item{
			Id:          insertUser.Id,
			Username:    insertUser.Username,
			PhoneNumber: insertUser.PhoneNumber,
			Email:       insertUser.Email,
			Name:        insertUser.Name,
			Avatar:      insertUser.Avatar,
			State:       insertUser.State,
			CreatedId:   insertUser.CreatedId,
			UpdatedId:   insertUser.UpdatedId,
			CreatedAt:   insertUser.CreatedAt.Format(global.TimeFormat),
			UpdatedAT:   insertUser.CreatedAt.Format(global.TimeFormat),
		},
	}, nil
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

func md5Password(password string) string {
	m := md5.New()
	m.Write([]byte(password))
	return fmt.Sprintf("%x", m.Sum(nil))
}

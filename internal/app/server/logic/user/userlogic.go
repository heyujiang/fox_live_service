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
		NiceName    string `json:"nice_name"`
		Avatar      string `json:"avatar"`
		State       int    `json:"state"`
		CreatedAt   string `json:"create_at"`
		UpdatedAT   string `json:"update_at"`
	}

	ReqCreateUser struct {
		Username    string `json:"username"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}

	RespCreateUser struct{}

	ReqUpdateUser struct {
		ReqUriUpdateUser
		ReqBodyUpdateUser
	}

	ReqUriUpdateUser struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateUser struct {
		Email    string `json:"email"`
		NickName string `json:"nick_name"`
		Avatar   string `json:"avatar"`
	}

	RespUpdateUser struct {
	}

	ReqUserInfo struct {
		Id int `uri:"id"`
	}

	RespUserInfo struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		NiceName    string `json:"nice_name"`
		Avatar      string `json:"avatar"`
		State       int    `json:"state"`
		CreatedAt   string `json:"create_at"`
		UpdatedAT   string `json:"update_at"`
	}

	ReqEnableUser struct {
		Id int `uri:"id"`
	}

	RespEnableUser struct {
	}

	ReqDisableUser struct {
		Id int `uri:"id"`
	}

	RespDisableUser struct {
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

	return &RespCreateUser{}, nil
}

// Update 修改用户信息
func (b *bisLogic) Update(req *ReqUpdateUser, uid int) (*RespUpdateUser, error) {
	_, err := model.UserModel.Find(req.Id)
	if err != nil {
		slog.Error("update user get user error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "用户不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户错误")
	}

	if err := model.UserModel.Update(&model.User{
		Id:        req.Id,
		Email:     req.Email,
		NickName:  req.NickName,
		Avatar:    req.Avatar,
		UpdatedId: uid,
	}); err != nil {
		slog.Error("update user error ", "id", req.Id, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改用户信息错误")
	}

	return &RespUpdateUser{}, nil
}

// Info 用户信息
func (b *bisLogic) Info(req *ReqUserInfo) (*RespUserInfo, error) {
	user, err := model.UserModel.Find(req.Id)
	if err != nil {
		slog.Error("update user get user error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "用户不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户错误")
	}

	return &RespUserInfo{
		Id:          user.Id,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Name:        user.Email,
		NiceName:    user.NickName,
		Avatar:      user.Avatar,
		State:       user.State,
		CreatedAt:   user.CreatedAt.Format(global.TimeFormat),
		UpdatedAT:   user.UpdatedAt.Format(global.TimeFormat),
	}, nil
}

func (b *bisLogic) List(req *ReqUserList) (*RespUserList, error) {
	logic.VerifyReqPage(&req.ReqPage)
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
			NiceName:    user.NickName,
			Avatar:      user.Avatar,
			State:       user.State,
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
	cond := &model.UserCond{}
	if req.Id > 0 {
		cond.Id = req.Id
	}

	if req.PhoneNumber != "" {
		cond.PhoneNumber = req.PhoneNumber
	}

	if req.Name != "" {
		cond.Name = req.Name
	}

	if _, ok := model.UserStatusDesc[req.State]; ok {
		cond.State = req.State
	}

	return cond
}

func md5Password(password string) string {
	m := md5.New()
	m.Write([]byte(password))
	return fmt.Sprintf("%x", m.Sum(nil))
}

func (b *bisLogic) Enable(req *ReqEnableUser, uid int) (*RespEnableUser, error) {
	_, err := model.UserModel.Find(req.Id)
	if err != nil {
		slog.Error("update user get user error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "用户不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户错误")
	}

	if err := model.UserModel.UpdateState(req.Id, model.UserStatusEnable, uid); err != nil {
		slog.Error("update user error ", "id", req.Id, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "启用用户错误")
	}

	return &RespEnableUser{}, nil
}

func (b *bisLogic) Disable(req *ReqDisableUser, uid int) (*RespDisableUser, error) {
	_, err := model.UserModel.Find(req.Id)
	if err != nil {
		slog.Error("update user get user error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "用户不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户错误")
	}

	if err := model.UserModel.UpdateState(req.Id, model.UserStatusDisable, uid); err != nil {
		slog.Error("update user error ", "id", req.Id, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "禁用用户错误")
	}

	return &RespDisableUser{}, nil
}

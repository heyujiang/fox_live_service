package user

import (
	"crypto/md5"
	"errors"
	"fmt"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/logic"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
	"strings"
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
		PhoneNumber string `form:"phoneNumber"`
		Username    string `form:"username"`
		State       int    `form:"state" binding:"omitempty,oneof=1 2"`
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
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		NickName    string `json:"nickName"`
		Avatar      string `json:"avatar"`
		State       int    `json:"state"`
		RoleIds     []int  `json:"roleIds"`
		Roles       string `json:"roles"`
		DeptId      int    `json:"deptId"`
		Dept        string `json:"dept"`
		Job         string `json:"job"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}

	ReqCreateUser struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required"`
		Email       string `json:"email" binding:"omitempty,email"`
		Name        string `json:"name"`
		NickName    string `json:"nickName"`
		Avatar      string `json:"avatar"`
		RoleIds     []int  `json:"roleIds" binding:"required"`
		DeptId      int    `json:"deptId" binding:"required"`
		Job         string `json:"job"`
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
		Email    string `json:"email" binding:"omitempty,email"`
		Name     string `json:"name"`
		NickName string `json:"nickName"`
		Avatar   string `json:"avatar"`
		RoleIds  []int  `json:"roleIds"`
		DeptId   int    `json:"deptId"`
		Job      string `json:"job"`
	}

	RespUpdateUser struct {
	}

	ReqUserInfo struct {
		Id int `uri:"id"`
	}

	RespUserInfo struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		NickName    string `json:"nickName"`
		Avatar      string `json:"avatar"`
		State       int    `json:"state"`
		RoleIds     []int  `json:"roleIds"`
		Roles       string `json:"roles"`
		DeptId      int    `json:"deptId"`
		Dept        string `json:"dept"`
		Job         string `json:"job"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}
	ReqDeleteUser struct {
		Id int `uri:"id"`
	}

	RespDeleteUser struct {
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

	RespOptionItem struct {
		Id       int    `json:"value"`
		Username string `json:"label"`
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
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Name:        req.Name,
		NickName:    req.NickName,
		Avatar:      req.Avatar,
		RoleIds:     strings.Join(cast.ToStringSlice(req.RoleIds), ","),
		DeptId:      req.DeptId,
		Job:         req.Job,
		IsSystem:    model.NonSystemUser,
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
		Name:      req.Name,
		Email:     req.Email,
		NickName:  req.NickName,
		Avatar:    req.Avatar,
		RoleIds:   strings.Join(cast.ToStringSlice(req.RoleIds), ","),
		DeptId:    req.DeptId,
		Job:       req.Job,
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

	dept, err := model.DeptModel.Find(user.DeptId)
	if err != nil && !errors.Is(err, model.ErrNotRecord) {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户部门错误")
	}

	roleIds := cast.ToIntSlice(strings.Split(user.RoleIds, ","))
	roles, err := model.RoleModel.SelectByIds(roleIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户角色错误")
	}

	roleTitles := make([]string, 0, len(roles))
	for _, role := range roles {
		roleTitles = append(roleTitles, role.Title)
	}

	return &RespUserInfo{
		Id:          user.Id,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Name:        user.Email,
		NickName:    user.NickName,
		Avatar:      user.Avatar,
		State:       user.State,
		RoleIds:     roleIds,
		Roles:       strings.Join(roleTitles, ","),
		DeptId:      user.DeptId,
		Dept:        dept.Title,
		Job:         user.Job,
		CreatedAt:   user.CreatedAt.Format(global.TimeFormat),
		UpdatedAt:   user.UpdatedAt.Format(global.TimeFormat),
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

	deptIds := make([]int, 0, len(users))
	roleIdStr := "0"
	for _, user := range users {
		deptIds = append(deptIds, user.DeptId)
		if len(user.RoleIds) > 0 {
			roleIdStr = fmt.Sprintf("%s,%s", roleIdStr, user.RoleIds)
		}
	}
	roleIds := cast.ToIntSlice(strings.Split(roleIdStr, ","))
	roles, err := model.RoleModel.SelectByIds(roleIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取用户角色出错")
	}
	roleMap := make(map[int]string, len(roles))
	for _, role := range roles {
		roleMap[role.Id] = role.Title
	}

	depts, err := model.DeptModel.SelectByIds(deptIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取用户部门出错")
	}
	deptMap := make(map[int]string, len(depts))
	for _, dept := range depts {
		deptMap[dept.Id] = dept.Title
	}

	items := make([]*Item, 0, len(users))
	for _, user := range users {
		userRoleIds := cast.ToIntSlice(strings.Split(user.RoleIds, ","))
		userRoles := make([]string, 0, len(userRoleIds))
		for _, roleId := range userRoleIds {
			if role, ok := roleMap[roleId]; ok {
				userRoles = append(userRoles, role)
			}
		}

		items = append(items, &Item{
			Id:          user.Id,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Name:        user.Name,
			NickName:    user.NickName,
			Avatar:      user.Avatar,
			State:       user.State,
			RoleIds:     userRoleIds,
			Roles:       strings.Join(userRoles, ","),
			Dept:        deptMap[user.DeptId],
			DeptId:      user.DeptId,
			Job:         user.Job,
			CreatedAt:   user.CreatedAt.Format(global.TimeFormat),
			UpdatedAt:   user.CreatedAt.Format(global.TimeFormat),
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

	if req.Username != "" {
		cond.Username = req.Username
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

func (b *bisLogic) Delete(req *ReqDeleteUser, uid int) (*RespDeleteUser, error) {
	_, err := model.UserModel.Find(req.Id)
	if err != nil {
		slog.Error("delete user get user error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "用户不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户错误")
	}

	if err := model.UserModel.Delete(req.Id); err != nil {
		slog.Error("delete user error ", "id", req.Id, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除用户错误")
	}

	return &RespDeleteUser{}, nil
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

func (b *bisLogic) Options() ([]*RespOptionItem, error) {
	users, err := model.UserModel.SelectByEnable()
	if err != nil {
		slog.Error("get user options error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户选项列表错误")
	}

	res := make([]*RespOptionItem, 0, len(users))
	for _, v := range users {
		res = append(res, &RespOptionItem{
			Id:       v.Id,
			Username: v.Username,
		})
	}

	return res, nil
}

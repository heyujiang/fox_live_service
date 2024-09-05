package user

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
	"strings"
)

var AccountLogic = newAccountLogic()

type (
	accountLogic struct{}

	RespAccountUserInfo struct {
		UserId      int    `json:"userId"`
		Username    string `json:"username"`
		Name        string `json:"name"`
		Nickname    string `json:"nickname"`
		Avatar      string `json:"avatar"`
		Dept        string `json:"dept"`
		Role        string `json:"role"`
		Email       string `json:"email"`
		Job         string `json:"job"`
		PhoneNumber string `json:"phoneNumber"`
		CreatedAt   string `json:"createdAt"`
	}

	ReqUpdateAvatar struct {
		Avatar string `json:"avatar"`
	}
	RespUpdateAvatar struct {
	}

	ReqUpdateBasic struct {
		Email string `json:"email" binding:"omitempty,email"`
	}
	RespUpdateBasic struct {
	}

	ReqUpdatePassword struct {
		OldPassword string `json:"oldPassword"`
		Password    string `json:"password"`
	}
	RespUpdatePassword struct {
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
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "账户不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取账户信息出错")
	}

	roleIds := cast.ToIntSlice(strings.Split(user.RoleIds, ","))
	roles, err := model.RoleModel.SelectByIds(roleIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取账户信息出错")
	}
	roleNames := make([]string, 0, len(roles))
	for _, v := range roles {
		roleNames = append(roleNames, v.Title)
	}

	dept, err := model.DeptModel.Find(user.DeptId)
	if err != nil {
		slog.Error("get account user info error ： ", "uid", uid, "err", err)
		if !errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "获取账户信息出错")
		}
	}

	return &RespAccountUserInfo{
		UserId:      user.Id,
		Username:    user.Username,
		Name:        user.Name,
		Nickname:    user.NickName,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Dept:        dept.Title,
		Role:        strings.Join(roleNames, "，"),
		Job:         user.Job,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt.Format(global.TimeFormat),
	}, nil

}

func (a *accountLogic) UpdateAvatar(req *ReqUpdateAvatar, uid int) (*RespUpdateAvatar, error) {
	if err := model.UserModel.UpdateAvatar(uid, req.Avatar); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改头像出错")
	}
	return &RespUpdateAvatar{}, nil
}

func (a *accountLogic) UpdateBasic(req *ReqUpdateBasic, uid int) (*RespUpdateBasic, error) {
	user, err := model.UserModel.Find(uid)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改异常！")
	}
	user.Email = req.Email
	user.UpdatedId = uid
	if err := model.UserModel.UpdateBasic(user); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改头像出错")
	}
	return &RespUpdateBasic{}, nil
}

func (a *accountLogic) UpdatePassword(req *ReqUpdatePassword, uid int) (*RespUpdatePassword, error) {
	user, err := model.UserModel.Find(uid)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改异常！")
	}
	if user.Password != req.OldPassword {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "原密码错误")
	}
	if err := model.UserModel.UpdatePassword(uid, req.Password); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改密码出错")
	}
	return &RespUpdatePassword{}, nil
}

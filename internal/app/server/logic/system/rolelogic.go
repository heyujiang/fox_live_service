package system

import (
	"errors"
	"fmt"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"github.com/spf13/cast"
	"strings"
	"time"
)

var (
	RoleLogic = newRoleLogic()
)

type (
	roleLogic struct {
	}

	ReqCreateRole struct {
		Title   string `json:"title"`
		Pid     int    `json:"pid"`
		Remark  string `json:"remark"`
		RuleIds []int  `json:"RuleIds"`
	}

	RespCreateRole struct{}

	ReqUpdateRole struct {
		ReqUpdateRoleUri
		ReqUpdateRoleBody
	}

	ReqUpdateRoleUri struct {
		Id int `uri:"id"`
	}

	ReqUpdateRoleBody struct {
		Title   string `json:"title"`
		Pid     int    `json:"pid"`
		Remark  string `json:"remark"`
		RuleIds []int  `json:"ruleIds"`
	}

	RespUpdateRole struct{}

	ReqDeleteRole struct {
		Id int `uri:"id"`
	}
	RespDeleteRole struct{}

	ReqUpdateRoleStatus struct {
		ReqUpdateRoleUri
		ReqUpdateRoleStatusBody
	}

	ReqUpdateRoleStatusBody struct {
		Status int `json:"status"`
	}

	RespUpdateRoleStatus struct {
	}

	RespRoleListItem struct {
		Id        int                 `json:"id"`
		Title     string              `json:"title"`
		Pid       int                 `json:"pid"`
		Status    int                 `json:"status"`
		Remark    string              `json:"remark"`
		RuleIds   []int               `json:"ruleIds"`
		CreatedId int                 `json:"createdId"`
		UpdatedId int                 `json:"updatedId"`
		CreatedAt string              `json:"createdAt"`
		UpdatedAt string              `json:"updatedAt"`
		Children  []*RespRoleListItem `json:"children"`
	}

	RespRoleParentItem struct {
		Id       int                   `json:"id"`
		Title    string                `json:"title"`
		Children []*RespRoleParentItem `json:"children"`
	}

	ReqGetRoleRules struct {
		Id int `uri:"id"`
	}
)

func newRoleLogic() *roleLogic {
	return &roleLogic{}
}

func (r *roleLogic) Create(req *ReqCreateRole, uid int) (*RespCreateRole, error) {
	if err := model.RoleModel.Create(&model.Role{
		Title:     req.Title,
		Pid:       req.Pid,
		Status:    model.RoleStatusEnable,
		Remark:    req.Remark,
		RuleIds:   strings.Join(cast.ToStringSlice(req.RuleIds), ","),
		CreatedId: uid,
		UpdatedId: uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建角色出错")
	}
	return &RespCreateRole{}, nil
}

func (r *roleLogic) Update(req *ReqUpdateRole, uid int) (*RespUpdateRole, error) {
	_, err := model.RoleModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "角色不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询角色信息出错")
	}

	if err := model.RoleModel.Update(&model.Role{
		Id:        req.Id,
		Title:     req.Title,
		Pid:       req.Pid,
		Remark:    req.Remark,
		RuleIds:   strings.Join(cast.ToStringSlice(req.RuleIds), ","),
		CreatedId: 0,
		UpdatedId: uid,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "编辑角色出错")
	}

	return &RespUpdateRole{}, nil
}

func (r *roleLogic) Delete(req *ReqDeleteRole) (*RespDeleteRole, error) {
	_, err := model.RoleModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "角色不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询角色信息出错")
	}

	if err := model.RoleModel.Delete(req.Id); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除角色出错")
	}
	return &RespDeleteRole{}, nil
}

func (r *roleLogic) List() ([]*RespRoleListItem, error) {
	roles, err := model.RoleModel.Select()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询角色出错")
	}

	roleMap := make(map[int][]*RespRoleListItem)
	roleMap[0] = []*RespRoleListItem{}
	for _, role := range roles {
		fmt.Println(fmt.Sprintf("%+v", role))
		if _, ok := roleMap[role.Pid]; !ok {
			roleMap[role.Pid] = make([]*RespRoleListItem, 0)
		}

		roleMap[role.Pid] = append(roleMap[role.Pid], &RespRoleListItem{
			Id:        role.Id,
			Title:     role.Title,
			Pid:       role.Pid,
			Remark:    role.Remark,
			Status:    role.Status,
			RuleIds:   cast.ToIntSlice(strings.Split(role.RuleIds, ",")),
			CreatedId: role.CreatedId,
			UpdatedId: role.UpdatedId,
			CreatedAt: role.CreatedAt.Format(global.TimeFormat),
			UpdatedAt: role.UpdatedAt.Format(global.TimeFormat),
		})
	}

	for _, roles := range roleMap {
		for _, role := range roles {
			if _, ok := roleMap[role.Id]; !ok {
				role.Children = make([]*RespRoleListItem, 0)
			} else {
				role.Children = roleMap[role.Id]
			}
		}
	}

	return roleMap[0], nil
}

func (r *roleLogic) UpdateStatus(req *ReqUpdateRoleStatus, uid int) (*RespUpdateRoleStatus, error) {
	_, err := model.RoleModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "角色不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询角色信息出错")
	}

	if err := model.RoleModel.UpdateStatus(req.Id, req.Status, uid); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改角色状态出错")
	}
	return &RespUpdateRoleStatus{}, nil
}

func (r *roleLogic) Parents() ([]*RespRoleParentItem, error) {
	roles, err := model.RoleModel.SelectEnable()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询菜单出错")
	}

	roleMap := make(map[int][]*RespRoleParentItem)
	roleMap[0] = []*RespRoleParentItem{}
	for _, role := range roles {
		fmt.Println(fmt.Sprintf("%+v", role))
		if _, ok := roleMap[role.Pid]; !ok {
			roleMap[role.Pid] = make([]*RespRoleParentItem, 0)
		}
		roleMap[role.Pid] = append(roleMap[role.Pid], &RespRoleParentItem{
			Id:    role.Id,
			Title: role.Title,
		})
	}

	for _, roles := range roleMap {
		for _, role := range roles {
			if _, ok := roleMap[role.Id]; !ok {
				role.Children = make([]*RespRoleParentItem, 0)
			} else {
				role.Children = roleMap[role.Id]
			}
		}
	}

	return roleMap[0], nil
}

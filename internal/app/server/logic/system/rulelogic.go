package system

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"github.com/spf13/cast"
	"strings"
)

var (
	RuleLogic = newRuleLogic()
)

type (
	ruleLogic struct {
	}

	ReqCreateRule struct {
		ActiveMenu         int    `json:"activeMenu"`         // 高亮菜单 （高亮设置的菜单项）"activemenu": 0,
		Component          string `json:"component"`          // 组件路径 "component": "system/usersetting/index.vue",
		HideChildrenInMenu int    `json:"hideChildrenInMenu"` // 显示单项 （hideChildrenInMenu强制在左侧菜单中显示单项） "hidechildreninmenu": 0,
		HideInMenu         int    `json:"hideInMenu"`         // 菜单隐藏（是否在左侧菜单中隐藏该项）"hideinmenu": 1,
		Icon               string `json:"icon"`               // 图标"icon": "icon-user",
		IsExt              int    `json:"isExt"`              // 是否外链"isext": 0,
		Keepalive          int    `json:"keepalive"`          // 是否缓存（在页签模式生效,页面name和路由name保持一致）"keepalive": 0,
		Locale             string `json:"locale"`             // 多语言标识 "locale": "",
		NoAffix            int    `json:"noAffix"`            // 不添加tab中 （如果设置为true标签将不会添加到tab-bar中） "noaffix": 0,
		OnlyPage           int    `json:"onlyPage"`           // 独立页面 （不需layout和登录，如登录页、数据大屏） "onlypage": 0,
		Order              int    `json:"order"`              // 排序"order": 1,
		Permission         string `json:"permission"`         // 权限标识 "permission": "",
		Pid                int    `json:"pid"`                // 上级菜单"pid": 0,
		Redirect           string `json:"redirect"`           // 重定向（菜单目录需要）"redirect": "",
		RequiresAuth       int    `json:"requiresAuth"`       // 是否登录鉴权（是否需要登录鉴权） "requiresauth": 1,
		RouteName          string `json:"routeName"`          // 路由名称 "routename": "usersetting",
		RoutePath          string `json:"routePath"`          // 路由地址 "routepath": "/usersetting",
		Spacer             string `json:"spacer"`             // "spacer": "",
		Status             int    `json:"status"`             // 状态 "status": 0,
		Title              string `json:"title"`              // 菜单名称 "title": "个人中心",
		Type               int    `json:"type"`               // 菜单类型 （菜单目录 菜单 按钮）"type": 1,
	}

	RespCreateRule struct{}

	ReqUpdateRule struct {
		ReqUpdateRuleUri
		ReqUpdateRuleBody
	}

	ReqUpdateRuleUri struct {
		Id int `uri:"id"`
	}

	ReqUpdateRuleBody struct {
		ActiveMenu         int    `json:"activeMenu"`         // 高亮菜单 （高亮设置的菜单项）"activemenu": 0,
		Component          string `json:"component"`          // 组件路径 "component": "system/usersetting/index.vue",
		HideChildrenInMenu int    `json:"hideChildrenInMenu"` // 显示单项 （hideChildrenInMenu强制在左侧菜单中显示单项） "hidechildreninmenu": 0,
		HideInMenu         int    `json:"hideInMenu"`         // 菜单隐藏（是否在左侧菜单中隐藏该项）"hideinmenu": 1,
		Icon               string `json:"icon"`               // 图标"icon": "icon-user",
		IsExt              int    `json:"isExt"`              // 是否外链"isext": 0,
		Keepalive          int    `json:"keepalive"`          // 是否缓存（在页签模式生效,页面name和路由name保持一致）"keepalive": 0,
		Locale             string `json:"locale"`             // 多语言标识 "locale": "",
		NoAffix            int    `json:"noAffix"`            // 不添加tab中 （如果设置为true标签将不会添加到tab-bar中） "noaffix": 0,
		OnlyPage           int    `json:"onlyPage"`           // 独立页面 （不需layout和登录，如登录页、数据大屏） "onlypage": 0,
		Order              int    `json:"order"`              // 排序"order": 1,
		Permission         string `json:"permission"`         // 权限标识 "permission": "",
		Pid                int    `json:"pid"`                // 上级菜单"pid": 0,
		Redirect           string `json:"redirect"`           // 重定向（菜单目录需要）"redirect": "",
		RequiresAuth       int    `json:"requiresAuth"`       // 是否登录鉴权（是否需要登录鉴权） "requiresauth": 1,
		RouteName          string `json:"routeName"`          // 路由名称 "routename": "usersetting",
		RoutePath          string `json:"routePath"`          // 路由地址 "routepath": "/usersetting",
		Spacer             string `json:"spacer"`             // "spacer": "",
		Status             int    `json:"status"`             // 状态 "status": 0,
		Title              string `json:"title"`              // 菜单名称 "title": "个人中心",
		Type               int    `json:"type"`               // 菜单类型 （菜单目录 菜单 按钮）"type": 1,
	}

	RespUpdateRule struct{}

	ReqDeleteRule struct {
		Id int `uri:"id"`
	}
	RespDeleteRule struct{}

	ReqUpdateRuleStatus struct {
		ReqUpdateRuleUri
		ReqUpdateRuleStatusBody
	}

	ReqUpdateRuleStatusBody struct {
		Status int `json:"status"`
	}

	RespUpdateRuleStatus struct {
	}

	RespRuleListItem struct {
		Id                 int                 `json:"id"`
		ActiveMenu         int                 `json:"activeMenu"`         // 高亮菜单 （高亮设置的菜单项）"activemenu": 0,
		Component          string              `json:"component"`          // 组件路径 "component": "system/usersetting/index.vue",
		HideChildrenInMenu int                 `json:"hideChildrenInMenu"` // 显示单项 （hideChildrenInMenu强制在左侧菜单中显示单项） "hidechildreninmenu": 0,
		HideInMenu         int                 `json:"hideInMenu"`         // 菜单隐藏（是否在左侧菜单中隐藏该项）"hideinmenu": 1,
		Icon               string              `json:"icon"`               // 图标"icon": "icon-user",
		IsExt              int                 `json:"isExt"`              // 是否外链"isext": 0,
		Keepalive          int                 `json:"keepalive"`          // 是否缓存（在页签模式生效,页面name和路由name保持一致）"keepalive": 0,
		Locale             string              `json:"locale"`             // 多语言标识 "locale": "",
		NoAffix            int                 `json:"noAffix"`            // 不添加tab中 （如果设置为true标签将不会添加到tab-bar中） "noaffix": 0,
		OnlyPage           int                 `json:"onlyPage"`           // 独立页面 （不需layout和登录，如登录页、数据大屏） "onlypage": 0,
		Order              int                 `json:"order"`              // 排序"order": 1,
		Permission         string              `json:"permission"`         // 权限标识 "permission": "",
		Pid                int                 `json:"pid"`                // 上级菜单"pid": 0,
		Redirect           string              `json:"redirect"`           // 重定向（菜单目录需要）"redirect": "",
		RequiresAuth       int                 `json:"requiresAuth"`       // 是否登录鉴权（是否需要登录鉴权） "requiresauth": 1,
		RouteName          string              `json:"routeName"`          // 路由名称 "routename": "usersetting",
		RoutePath          string              `json:"routePath"`          // 路由地址 "routepath": "/usersetting",
		Spacer             string              `json:"spacer"`             // "spacer": "",
		Status             int                 `json:"status"`             // 状态 "status": 0,
		Title              string              `json:"title"`              // 菜单名称 "title": "个人中心",
		Type               int                 `json:"type"`               // 菜单类型 （菜单目录 菜单 按钮）"type": 1,
		CreatedId          int                 `json:"createdId"`
		UpdatedId          int                 `json:"updatedId"`
		CreatedAt          string              `json:"createdAt"`
		UpdatedAt          string              `json:"updatedAt"`
		Children           []*RespRuleListItem `json:"children"`
	}

	RespRuleParentItem struct {
		Id       int                   `json:"id"`
		Title    string                `json:"title"`
		Children []*RespRuleParentItem `json:"children"`
	}
)

func newRuleLogic() *ruleLogic {
	return &ruleLogic{}
}

func (r *ruleLogic) Create(req *ReqCreateRule, uid int) (*RespCreateRule, error) {
	if err := model.RuleModel.Create(&model.Rule{
		Title:              req.Title,
		RouteName:          req.RouteName,
		RoutePath:          req.RoutePath,
		Component:          req.Component,
		Redirect:           req.Redirect,
		Locale:             req.Locale,
		Icon:               req.Icon,
		Permission:         req.Permission,
		Spacer:             req.Spacer,
		HideInMenu:         req.HideInMenu,
		IsExt:              req.IsExt,
		NoAffix:            req.NoAffix,
		Keepalive:          req.Keepalive,
		RequiresAuth:       req.RequiresAuth,
		OnlyPage:           req.OnlyPage,
		ActiveMenu:         req.ActiveMenu,
		HideChildrenInMenu: req.HideChildrenInMenu,
		Status:             model.RuleStatusEnable,
		Order:              req.Order,
		Type:               req.Type,
		Pid:                req.Pid,
		CreatedId:          uid,
		UpdatedId:          uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建菜单出错")
	}
	return &RespCreateRule{}, nil
}

func (r *ruleLogic) Update(req *ReqUpdateRule, uid int) (*RespUpdateRule, error) {
	_, err := model.RuleModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "菜单不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询菜单信息出错")
	}
	if err := model.RuleModel.Update(&model.Rule{
		Id:                 req.Id,
		Title:              req.Title,
		RouteName:          req.RouteName,
		RoutePath:          req.RoutePath,
		Component:          req.Component,
		Redirect:           req.Redirect,
		Locale:             req.Locale,
		Icon:               req.Icon,
		Permission:         req.Permission,
		Spacer:             req.Spacer,
		HideInMenu:         req.HideInMenu,
		IsExt:              req.IsExt,
		NoAffix:            req.NoAffix,
		Keepalive:          req.Keepalive,
		RequiresAuth:       req.RequiresAuth,
		OnlyPage:           req.OnlyPage,
		ActiveMenu:         req.ActiveMenu,
		HideChildrenInMenu: req.HideChildrenInMenu,
		Order:              req.Order,
		Type:               req.Type,
		Pid:                req.Pid,
		UpdatedId:          uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "编辑菜单出错")
	}

	return &RespUpdateRule{}, nil
}

func (r *ruleLogic) Delete(req *ReqDeleteRule) (*RespDeleteRule, error) {
	_, err := model.RuleModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "菜单不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询菜单信息出错")
	}

	if err := model.RuleModel.Delete(req.Id); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除菜单出错")
	}
	return &RespDeleteRule{}, nil
}

func (r *ruleLogic) List() ([]*RespRuleListItem, error) {
	rules, err := model.RuleModel.Select()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询菜单出错")
	}

	ruleMap := make(map[int][]*RespRuleListItem)
	ruleMap[0] = []*RespRuleListItem{}
	for _, rule := range rules {
		if _, ok := ruleMap[rule.Pid]; !ok {
			ruleMap[rule.Pid] = make([]*RespRuleListItem, 0)
		}
		ruleMap[rule.Pid] = append(ruleMap[rule.Pid], &RespRuleListItem{
			Id:                 rule.Id,
			ActiveMenu:         rule.ActiveMenu,
			Component:          rule.Component,
			HideChildrenInMenu: rule.HideChildrenInMenu,
			HideInMenu:         rule.HideInMenu,
			Icon:               rule.Icon,
			IsExt:              rule.IsExt,
			Keepalive:          rule.Keepalive,
			Locale:             rule.Locale,
			NoAffix:            rule.NoAffix,
			OnlyPage:           rule.OnlyPage,
			Order:              rule.Order,
			Permission:         rule.Permission,
			Pid:                rule.Pid,
			Redirect:           rule.Redirect,
			RequiresAuth:       rule.RequiresAuth,
			RouteName:          rule.RouteName,
			RoutePath:          rule.RoutePath,
			Spacer:             rule.Spacer,
			Status:             rule.Status,
			Title:              rule.Title,
			Type:               rule.Type,
			CreatedId:          rule.CreatedId,
			UpdatedId:          rule.UpdatedId,
			CreatedAt:          rule.CreatedAt.Format(global.TimeFormat),
			UpdatedAt:          rule.UpdatedAt.Format(global.TimeFormat),
		})
	}

	for _, rules := range ruleMap {
		for _, rule := range rules {
			if _, ok := ruleMap[rule.Id]; !ok {
				rule.Children = make([]*RespRuleListItem, 0)
			} else {
				rule.Children = ruleMap[rule.Id]
			}
		}
	}

	return ruleMap[0], nil
}

func (r *ruleLogic) UpdateStatus(req *ReqUpdateRuleStatus, uid int) (*RespUpdateRuleStatus, error) {
	_, err := model.RuleModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "菜单不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询菜单信息出错")
	}

	if err := model.RuleModel.UpdateStatus(req.Id, req.Status, uid); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改菜单状态出错")
	}
	return &RespUpdateRuleStatus{}, nil
}

func (r *ruleLogic) Parents() ([]*RespRuleParentItem, error) {
	rules, err := model.RuleModel.SelectDirAndMenu()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询菜单出错")
	}

	ruleMap := make(map[int][]*RespRuleParentItem)
	ruleMap[0] = []*RespRuleParentItem{}
	for _, rule := range rules {
		if _, ok := ruleMap[rule.Pid]; !ok {
			ruleMap[rule.Pid] = make([]*RespRuleParentItem, 0)
		}
		ruleMap[rule.Pid] = append(ruleMap[rule.Pid], &RespRuleParentItem{
			Id:    rule.Id,
			Title: rule.Title,
		})
	}

	for _, rules := range ruleMap {
		for _, rule := range rules {
			if _, ok := ruleMap[rule.Id]; !ok {
				rule.Children = make([]*RespRuleParentItem, 0)
			} else {
				rule.Children = ruleMap[rule.Id]
			}
		}
	}

	return ruleMap[0], nil
}

func (r *ruleLogic) GetRules(req *ReqGetRoleRules) ([]*RespRuleParentItem, error) {
	role, err := model.RoleModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "角色不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询角色信息出错")
	}

	rules := make([]*model.Rule, 0)
	rules, err = model.RuleModel.SelectEnable()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询菜单信息出错")
	}

	roleRuleIdMap := make(map[int]struct{}) // role 拥有的 rule 的 id 的map ， 用作验证 rule 是否属于当前角色
	if role.Id != 1 {
		roleIds := cast.ToIntSlice(strings.Split(role.RuleIds, ","))
		for _, roleId := range roleIds {
			roleRuleIdMap[roleId] = struct{}{}
		}

		ruleMapKId := make(map[int]*model.Rule, len(rules))
		for _, rule := range rules {
			ruleMapKId[rule.Id] = rule
		}

		for _, rule := range rules {
			if _, ok := roleRuleIdMap[rule.Id]; ok { // rule 是当前 role 的 rule
				//判断是否是顶级rule ，如否，判断其PID是否属于当前角色的rule，如否加入
				if rule.Pid != 0 {
					if _, pok := roleRuleIdMap[rule.Pid]; !pok {
						roleIds = append(roleIds, rule.Pid)
						roleRuleIdMap[rule.Pid] = struct{}{}
					}
				}
			}
		}

		newRules := make([]*model.Rule, 0)
		for _, rule := range rules {
			if _, ok := roleRuleIdMap[rule.Id]; ok {
				newRules = append(newRules, rule)
			}
		}
		rules = newRules
	}

	ruleMap := make(map[int][]*RespRuleParentItem)
	ruleMap[0] = []*RespRuleParentItem{}
	for _, rule := range rules {
		if _, ok := ruleMap[rule.Pid]; !ok {
			ruleMap[rule.Pid] = make([]*RespRuleParentItem, 0)
		}
		ruleMap[rule.Pid] = append(ruleMap[rule.Pid], &RespRuleParentItem{
			Id:    rule.Id,
			Title: rule.Title,
		})
	}

	for _, rules := range ruleMap {
		for _, rule := range rules {
			if _, ok := ruleMap[rule.Id]; !ok {
				rule.Children = make([]*RespRuleParentItem, 0)
			} else {
				rule.Children = ruleMap[rule.Id]
			}
		}
	}

	return ruleMap[0], nil
}

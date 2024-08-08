package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
	"time"
)

const (
	RuleTypeDir = iota + 1
	RuleTypeMenu
	RuleTypeBut
)

const (
	RuleStatusEnable = iota + 1
	RuleStatusDisable
)

const (
	insertRuleStr = "`title` , `route_name` , `route_path`, `component` , `redirect` , `locale` , `icon`, `permission`, `spacer` , `hide_in_menu` , `is_ext` , `no_affix` , `keepalive` , `requires_auth` , `only_page` , `active_menu`, `hide_children_in_menu` , `status` , `order` , `type` , `pid` , `created_id`, `updated_id` "
)

var (
	RuleModel = newRuleModel()
)

type (
	Rule struct {
		Id                 int       `db:"id"`
		Title              string    `db:"title"`                 // 菜单名称 "title": "个人中心",
		RouteName          string    `db:"route_name"`            // 路由名称 "routename": "usersetting",
		RoutePath          string    `db:"route_path"`            // 路由地址 "routepath": "/usersetting",
		Component          string    `db:"component"`             // 组件路径 "component": "system/usersetting/index.vue",
		Redirect           string    `db:"redirect"`              // 重定向（菜单目录需要）"redirect": "",
		Locale             string    `db:"locale"`                // 多语言标识 "locale": "",
		Icon               string    `db:"icon"`                  // 图标"icon": "icon-user",
		Permission         string    `db:"permission"`            // 权限标识 "permission": "",
		Spacer             string    `db:"spacer"`                // "spacer": "",
		HideInMenu         int       `db:"hide_in_menu"`          // 菜单隐藏（是否在左侧菜单中隐藏该项）"hideinmenu": 1,
		IsExt              int       `db:"is_ext"`                // 是否外链"isext": 0,
		NoAffix            int       `db:"no_affix"`              // 不添加tab中 （如果设置为true标签将不会添加到tab-bar中） "noaffix": 0,
		Keepalive          int       `db:"keepalive"`             // 是否缓存（在页签模式生效,页面name和路由name保持一致）"keepalive": 0,
		RequiresAuth       int       `db:"requires_auth"`         // 是否登录鉴权（是否需要登录鉴权） "requiresauth": 1,
		OnlyPage           int       `db:"only_page"`             // 独立页面 （不需layout和登录，如登录页、数据大屏） "onlypage": 0,
		ActiveMenu         int       `db:"active_menu"`           // 高亮菜单 （高亮设置的菜单项）"activemenu": 0,
		HideChildrenInMenu int       `db:"hide_children_in_menu"` // 显示单项 （hideChildrenInMenu强制在左侧菜单中显示单项） "hidechildreninmenu": 0,
		Status             int       `db:"status"`                // 状态 "status": 0,
		Order              int       `db:"order"`                 // 排序"order": 1,
		Type               int       `db:"type"`                  // 菜单类型 （菜单目录 菜单 按钮）"type": 1,
		Pid                int       `db:"pid"`                   // 上级菜单"pid": 0,
		CreatedId          int       `db:"created_id"`
		UpdatedId          int       `db:"updated_id"`
		CreatedAt          time.Time `db:"created_at"`
		UpdatedAt          time.Time `db:"updated_at"`
	}

	ruleModel struct {
		table string
	}
)

func newRuleModel() *ruleModel {
	return &ruleModel{
		table: "rule",
	}
}

func (m *ruleModel) Create(rule *Rule) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", m.table, insertRuleStr)
	res, err := db.Exec(sqlStr, rule.Title, rule.RouteName, rule.RoutePath, rule.Component, rule.Redirect, rule.Locale,
		rule.Icon, rule.Permission, rule.Spacer, rule.HideInMenu, rule.IsExt, rule.NoAffix, rule.Keepalive,
		rule.RequiresAuth, rule.OnlyPage, rule.ActiveMenu, rule.HideChildrenInMenu, rule.Status, rule.Order,
		rule.Type, rule.Pid, rule.CreatedId, rule.UpdatedId)
	if err != nil {
		slog.Error("insert rule err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastId, _ := res.LastInsertId()
	rule.Id = cast.ToInt(lastId)
	return nil
}

func (m *ruleModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete rule err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *ruleModel) Update(rule *Rule) error {
	sqlStr := fmt.Sprintf("update %s set `title` = ?, `route_name` = ? , `route_path` = ?, `component` = ? , `redirect` = ? , `locale` = ? , `icon` = ?, `permission` = ?, `spacer` = ? , `hide_in_menu` = ? , `is_ext` = ? , `no_affix` = ? , `keepalive` = ? , `requires_auth` = ? , `only_page` = ? , `active_menu` = ?, `hide_children_in_menu` = ? , `order` = ? , `type` = ? , `pid` = ? , updated_id = ? where `id` = %d", m.table, rule.Id)
	_, err := db.Exec(sqlStr, rule.Title, rule.RouteName, rule.RoutePath, rule.Component, rule.Redirect, rule.Locale,
		rule.Icon, rule.Permission, rule.Spacer, rule.HideInMenu, rule.IsExt, rule.NoAffix, rule.Keepalive,
		rule.RequiresAuth, rule.OnlyPage, rule.ActiveMenu, rule.HideChildrenInMenu, rule.Order,
		rule.Type, rule.Pid, rule.UpdatedId)
	if err != nil {
		slog.Error("update rule err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *ruleModel) UpdateStatus(id, status, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `status` = ?, `updated_id` = ?  where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, status, uid)
	if err != nil {
		slog.Error("update rule status err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *ruleModel) Find(id int) (*Rule, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
	rule := new(Rule)
	if err := db.Get(rule, sqlStr, id); err != nil {
		slog.Error("find rule err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return rule, nil
}

func (m *ruleModel) Select() ([]*Rule, error) {
	sqlStr := fmt.Sprintf("select * from %s where 1 = 1 order by `order` asc , id asc", m.table)
	rules := make([]*Rule, 0)
	if err := db.Select(&rules, sqlStr); err != nil {
		slog.Error("select rule err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return rules, nil
}

func (m *ruleModel) SelectDirAndMenu() ([]*Rule, error) {
	sqlStr := fmt.Sprintf("select * from %s where type in (?,?) and status = ? order by `order` asc , id asc", m.table)
	rules := make([]*Rule, 0)
	if err := db.Select(&rules, sqlStr, RuleTypeDir, RuleTypeMenu, RuleStatusEnable); err != nil {
		slog.Error("select rule err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return rules, nil
}

func (m *ruleModel) SelectByIds(ids []int) ([]*Rule, error) {
	var rules []*Rule
	sqlStr := fmt.Sprintf("select * from %s where `status` = ? and id in (?) ", m.table)
	query1, args, err := sqlx.In(sqlStr, RuleStatusEnable, ids)
	if err != nil {
		slog.Error("batch select rule by ids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}
	slog.Info(query1, "args", args)
	if err := db.Select(&rules, query1, args...); err != nil {
		slog.Error("batch select rule by ids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}
	return rules, nil
}

func (m *ruleModel) SelectEnable() ([]*Rule, error) {
	sqlStr := fmt.Sprintf("select * from %s where status = ? order by `order` asc , id asc", m.table)
	rules := make([]*Rule, 0)
	if err := db.Select(&rules, sqlStr, RuleStatusEnable); err != nil {
		slog.Error("select rule err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return rules, nil
}

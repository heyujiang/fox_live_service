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
	RoleStatusEnable = iota + 1
	RoleStatusDisable
)

const (
	insertRoleStr = "`title` , `pid` , `status`, `remark`, `rule_ids` ,  `created_id`, `updated_id` "
)

var (
	RoleModel = newRoleModel()
)

type (
	Role struct {
		Id        int       `db:"id"`
		Title     string    `db:"title"`
		Pid       int       `db:"pid"`
		Status    int       `db:"status"`
		Remark    string    `db:"remark"`
		RuleIds   string    `db:"rule_ids"`
		CreatedId int       `db:"created_id"`
		UpdatedId int       `db:"updated_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	roleModel struct {
		table string
	}
)

func newRoleModel() *roleModel {
	return &roleModel{
		table: "role",
	}
}

func (m *roleModel) Create(role *Role) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?)", m.table, insertRoleStr)
	res, err := db.Exec(sqlStr, role.Title, role.Pid, role.Status, role.Remark, role.RuleIds, role.CreatedId, role.UpdatedId)
	if err != nil {
		slog.Error("insert role err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastId, _ := res.LastInsertId()
	role.Id = cast.ToInt(lastId)
	return nil
}

func (m *roleModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete role err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *roleModel) Update(role *Role) error {
	sqlStr := fmt.Sprintf("update %s set `title` = ?, `pid` = ? , `remark` = ?, `rule_ids` = ? ,updated_id = ? where `id` = %d", m.table, role.Id)
	_, err := db.Exec(sqlStr, role.Title, role.Pid, role.Remark, role.RuleIds, role.UpdatedId)
	if err != nil {
		slog.Error("update role err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *roleModel) UpdateStatus(id, status, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `status` = ?, `updated_id` = ?  where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, status, uid)
	if err != nil {
		slog.Error("update rule status err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *roleModel) Find(id int) (*Role, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
	role := new(Role)
	if err := db.Get(role, sqlStr, id); err != nil {
		slog.Error("find role err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return role, nil
}

func (m *roleModel) Select() ([]*Role, error) {
	sqlStr := fmt.Sprintf("select * from %s where 1 = 1 order by id asc", m.table)
	roles := make([]*Role, 0)
	if err := db.Select(&roles, sqlStr); err != nil {
		slog.Error("select role err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return roles, nil
}

func (m *roleModel) SelectEnable() ([]*Role, error) {
	sqlStr := fmt.Sprintf("select * from %s where status = ? order by id asc", m.table)
	roles := make([]*Role, 0)
	if err := db.Select(&roles, sqlStr, RuleStatusEnable); err != nil {
		slog.Error("select role err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return roles, nil
}

func (m *roleModel) SelectByIds(ids []int) ([]*Role, error) {
	var roles []*Role
	if len(ids) == 0 {
		return roles, nil
	}
	sqlStr := fmt.Sprintf("select * from %s where id in (?) order by id asc", m.table)
	query, args, err := sqlx.In(sqlStr, ids)

	if err != nil {
		slog.Error("batch select role by ids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}

	if err := db.Select(&roles, query, args...); err != nil {
		slog.Error("batch select role by uids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}
	return roles, nil
}

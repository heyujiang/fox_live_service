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
	DeptStatusEnable = iota + 1
	DeptStatusDisable
)

const (
	insertDeptStr = "`name` , `pid` , `status`, `remark`, `order` , `created_id`, `updated_id` "
)

var (
	DeptModel = newDeptModel()
)

type (
	Dept struct {
		Id        int       `db:"id"`
		Name      string    `db:"name"`
		Pid       int       `db:"pid"`
		Status    int       `db:"status"`
		Remark    string    `db:"remark"`
		Order     int       `db:"order"`
		CreatedId int       `db:"created_id"`
		UpdatedId int       `db:"updated_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	deptModel struct {
		table string
	}
)

func newDeptModel() *deptModel {
	return &deptModel{
		table: "dept",
	}
}

func (m *deptModel) Create(dept *Dept) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?)", m.table, insertDeptStr)
	res, err := db.Exec(sqlStr, dept.Name, dept.Pid, dept.Status, dept.Remark, dept.Order, dept.CreatedId, dept.UpdatedId)
	if err != nil {
		slog.Error("insert dept err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastId, _ := res.LastInsertId()
	dept.Id = cast.ToInt(lastId)
	return nil
}

func (m *deptModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete dept err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *deptModel) Update(dept *Dept) error {
	sqlStr := fmt.Sprintf("update %s set `name` = ?, `pid` = ? , `remark` = ?, `order` = ? ,updated_id = ? where `id` = %d", m.table, dept.Id)
	_, err := db.Exec(sqlStr, dept.Name, dept.Pid, dept.Remark, dept.Order, dept.UpdatedId)
	if err != nil {
		slog.Error("update dept err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *deptModel) UpdateStatus(id, status, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `status` = ?, `updated_id` = ?  where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, status, uid)
	if err != nil {
		slog.Error("update dept status err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *deptModel) Find(id int) (*Dept, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
	dept := &Dept{}
	if err := db.Get(dept, sqlStr, id); err != nil {
		slog.Error("find dept err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return dept, ErrNotRecord
		}
		return nil, err
	}
	return dept, nil
}

func (m *deptModel) Select() ([]*Dept, error) {
	sqlStr := fmt.Sprintf("select * from %s where 1 = 1 order by id asc", m.table)
	depts := make([]*Dept, 0)
	if err := db.Select(&depts, sqlStr); err != nil {
		slog.Error("select dept err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return depts, nil
}

func (m *deptModel) SelectEnable() ([]*Dept, error) {
	sqlStr := fmt.Sprintf("select * from %s where status = ? order by id asc", m.table)
	depts := make([]*Dept, 0)
	if err := db.Select(&depts, sqlStr, RuleStatusEnable); err != nil {
		slog.Error("select Dept err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return depts, nil
}

func (m *deptModel) SelectByIds(ids []int) ([]*Dept, error) {
	var depts []*Dept
	if len(ids) == 0 {
		return depts, nil
	}

	sqlStr := fmt.Sprintf("select * from %s where id in (?)order by id asc", m.table)
	query, args, err := sqlx.In(sqlStr, ids)
	if err != nil {
		slog.Error("batch select dept by ids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}

	if err := db.Select(&depts, query, args...); err != nil {
		slog.Error("batch select dept by uids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}
	return depts, nil
}

package model

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"time"
)

const (
	inertProjectAttachedStr = "`project_id`,`node_id`,`record_id`,`user_id`,`url`,`filename`,`mime`,`size`,`created_id`"
)

var (
	ProjectAttachedModel = newProjectAttachedModel()
)

type (
	ProjectAttached struct {
		Id        int       `db:"id"`
		ProjectId int       `db:"project_id"`
		NodeId    int       `db:"node_id"`
		RecordId  int       `db:"record_id"`
		UserId    int       `db:"user_id"`
		Url       string    `db:"url"`
		Filename  string    `db:"filename"`
		Mime      string    `db:"mime"`
		Size      int64     `db:"size"`
		CreatedId int       `db:"created_id"`
		CreatedAt time.Time `db:"created_at"`
	}

	projectAttachedModel struct {
		table string
	}

	ProjectAttachedCond struct {
		Id        int
		RecordId  int
		ProjectId int
		NodeId    int
		UserId    int
	}
)

func newProjectAttachedModel() *projectAttachedModel {
	return &projectAttachedModel{
		table: "project_attached",
	}
}

func (m *projectAttachedModel) Create(projectAttached *ProjectAttached) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?)", m.table, inertProjectAttachedStr)
	_, err := db.Exec(sqlStr, projectAttached.ProjectId, projectAttached.NodeId, projectAttached.RecordId, projectAttached.UserId, projectAttached.Url, projectAttached.Filename, projectAttached.Mime, projectAttached.Size, projectAttached.CreatedId)
	if err != nil {
		slog.Error("insert project attached err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectAttachedModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete project attached err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectAttachedModel) Find(id int) (*ProjectAttached, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
	projectAttached := new(ProjectAttached)
	if err := db.Get(projectAttached, sqlStr, id); err != nil {
		slog.Error("find project attached err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return projectAttached, nil
}

func (m *projectAttachedModel) GetProjectAttachedByCond(cond *ProjectAttachedCond, pageIndex, pageSize int) ([]*ProjectAttached, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	sqlCond, args := m.buildProjectAttachedCond(cond)
	sqlStr := fmt.Sprintf("select * from %s where 1 = 1 %s limit %d,%d", m.table, sqlCond, (pageIndex-1)*pageSize, pageSize)
	var projectAttacheds []*ProjectAttached
	if err := db.Select(&projectAttacheds, sqlStr, args...); err != nil {
		slog.Error("get project attached error ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return projectAttacheds, nil
}

func (m *projectAttachedModel) buildProjectAttachedCond(cond *ProjectAttachedCond) (sqlCond string, args []interface{}) {
	if cond == nil {
		return
	}

	if cond.Id > 0 {
		sqlCond += "and id = ?"
		args = append(args, cond.Id)
	}
	if cond.ProjectId > 0 {
		sqlCond += "and project_id = ?"
		args = append(args, cond.ProjectId)
	}
	if cond.RecordId > 0 {
		sqlCond += "and record_id = ?"
		args = append(args, cond.RecordId)
	}
	if cond.NodeId > 0 {
		sqlCond += "and node_id = ?"
		args = append(args, cond.NodeId)
	}
	if cond.UserId > 0 {
		sqlCond += "and user_id = ?"
		args = append(args, cond.UserId)
	}
	return
}

func (m *projectAttachedModel) GetAllByProjectId(projectId int) ([]*ProjectAttached, error) {
	sqlStr := fmt.Sprintf("select * from %s where project_id = ?", m.table)
	var projectAttacheds []*ProjectAttached
	if err := db.Select(&projectAttacheds, sqlStr, projectId); err != nil {
		slog.Error("get project attached error ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
		return nil, err
	}
	return projectAttacheds, nil
}

package model

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"time"
)

const (
	inertProjectRecordStr = "`project_id`,`node_id`,`user_id`,`overview`,`created_at`,`updated_at`"
)

var (
	ProjectRecordModel = newProjectRecordModel()
)

type (
	ProjectRecord struct {
		Id        int       `db:"id"`
		ProjectId string    `db:"project_id"`
		NodeId    int       `db:"node_id"`
		UserId    string    `db:"user_id"`
		Overview  int       `db:"overview"`
		CreatedId int       `db:"created_id"`
		UpdatedId int       `db:"updated_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	projectRecordModel struct {
		table string
	}

	ProjectRecordCond struct {
		Id        int
		ProjectId int
		NodeId    int
		UserId    int
	}
)

func newProjectRecordModel() *projectRecordModel {
	return &projectRecordModel{
		table: "project_record",
	}
}

func (m *projectRecordModel) Create(projectRecord *ProjectRecord) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?)", m.table, inertProjectRecordStr)
	_, err := db.Exec(sqlStr, projectRecord.ProjectId, projectRecord.NodeId, projectRecord.UserId, projectRecord.Overview, projectRecord.CreatedId, projectRecord.UpdatedId)
	if err != nil {
		slog.Error("insert project record err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectRecordModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete project record err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectRecordModel) Update(projectRecord *ProjectRecord) error {
	sqlStr := fmt.Sprintf("update %s set `overview` = ? , `updated_id`= ? where `id` = %d", m.table, projectRecord.Id)
	_, err := db.Exec(sqlStr, projectRecord.Overview, projectRecord.UpdatedId)
	if err != nil {
		slog.Error("update project record err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectRecordModel) Find(id int) (*ProjectRecord, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
	projectRecord := new(ProjectRecord)
	if err := db.Get(projectRecord, sqlStr, id); err != nil {
		slog.Error("find project record err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return projectRecord, nil
}

func (m *projectRecordModel) GetProjectRecordByCond(cond *ProjectRecordCond, pageIndex, pageSize int) ([]*ProjectRecord, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	sqlCond, args := m.buildProjectRecordCond(cond)
	sqlStr := fmt.Sprintf("select * from %s where 1 = 1 %s limit %d,%d", m.table, sqlCond, (pageIndex-1)*pageSize, pageSize)
	var projectRecords []*ProjectRecord
	if err := db.Select(&projectRecords, sqlStr, args...); err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return projectRecords, nil
}

func (m *projectRecordModel) buildProjectRecordCond(cond *ProjectRecordCond) (sqlCond string, args []interface{}) {
	if cond == nil {
		return
	}

	if cond.Id > 0 {
		sqlCond += "and id = ?"
		args = append(args, cond.Id)
	}
	if cond.ProjectId > 0 {
		sqlCond += "and project_id = ?"
		args = append(args, cond.Id)
	}
	if cond.NodeId > 0 {
		sqlCond += "and node_id = ?"
		args = append(args, cond.Id)
	}
	if cond.UserId > 0 {
		sqlCond += "and user_id = ?"
		args = append(args, cond.Id)
	}
	return
}
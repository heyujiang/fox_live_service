package model

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"time"
)

const (
	ProjectRecordStateFinished int = iota + 1
	ProjectRecordStateIng

	inertProjectRecordStr = "`project_id`,`project_name`,`node_id`,`node_name`,`user_id`,`username`,`overview`,`state`,`created_id`,`updated_id`"
)

var (
	ProjectRecordModel = newProjectRecordModel()
)

type (
	ProjectRecord struct {
		Id          int       `db:"id"`
		ProjectId   int       `db:"project_id"`
		ProjectName string    `db:"project_name"`
		NodeId      int       `db:"node_id"`
		NodeName    string    `db:"node_name"`
		UserId      int       `db:"user_id"`
		Username    string    `db:"username"`
		Overview    string    `db:"overview"`
		State       int       `db:"state"`
		CreatedId   int       `db:"created_id"`
		UpdatedId   int       `db:"updated_id"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
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
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?,?)", m.table, inertProjectRecordStr)
	_, err := db.Exec(sqlStr, projectRecord.ProjectId, projectRecord.ProjectName, projectRecord.NodeId, projectRecord.NodeName, projectRecord.UserId, projectRecord.Username, projectRecord.Overview, projectRecord.State, projectRecord.CreatedId, projectRecord.UpdatedId)
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

func (m *projectRecordModel) GetProjectRecordCountByCond(cond *ProjectRecordCond) (int, error) {
	sqlCond, args := m.buildProjectRecordCond(cond)
	sqlStr := fmt.Sprintf("select count(*) from %s where 1 = 1 %s order by created_at desc", m.table, sqlCond)
	var count int
	if err := db.Get(&count, sqlStr, args...); err != nil {
		slog.Error("get project record count error ", "sql", sqlStr, "err ", err.Error())
		return 0, err
	}
	return count, nil
}

func (m *projectRecordModel) GetProjectRecordByCond(cond *ProjectRecordCond, pageIndex, pageSize int) ([]*ProjectRecord, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	sqlCond, args := m.buildProjectRecordCond(cond)
	fmt.Println(fmt.Sprintf("sqlcond : %s , args : %+v", sqlCond, args))
	sqlStr := fmt.Sprintf("select * from %s where 1 = 1 %s limit %d,%d", m.table, sqlCond, (pageIndex-1)*pageSize, pageSize)
	fmt.Println(sqlStr)
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
	fmt.Println(fmt.Sprintf("%+v", cond))
	if cond.Id > 0 {
		sqlCond += " and id = ?"
		args = append(args, cond.Id)
	}
	if cond.ProjectId > 0 {
		sqlCond += " and project_id = ?"
		args = append(args, cond.ProjectId)
	}
	if cond.NodeId > 0 {
		sqlCond += " and node_id = ?"
		args = append(args, cond.NodeId)
	}
	if cond.UserId > 0 {
		sqlCond += " and user_id = ?"
		args = append(args, cond.UserId)
	}
	return
}

func (m *projectRecordModel) GetAllByProjectId(projectId int) ([]*ProjectRecord, error) {
	sqlStr := fmt.Sprintf("select * from %s where project_id = ?", m.table)
	var projectRecords []*ProjectRecord
	if err := db.Select(&projectRecords, sqlStr, projectId); err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
		return nil, err
	}
	return projectRecords, nil
}

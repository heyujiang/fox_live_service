package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"
	"strings"
	"time"
)

const (
	insertProjectPersonStr = "`project_id`,`user_id`,`type`,`name`,`phone_number`,`created_id`"
)

const (
	ProjectPersonTypeFirst  = iota + 1 // 第一负责人
	ProjectPersonTypeSecond            // 第二负责人
	ProjectPersonTypeCommon            // 项目成员
)

var (
	ProjectPersonModel = newProjectPersonModel()
)

type (
	ProjectPerson struct {
		Id          int       `db:"id"`
		ProjectId   int       `db:"project_id"`
		UserId      int       `db:"user_id"`
		Type        int       `db:"type"`
		Name        string    `db:"name"`
		PhoneNumber string    `db:"phone_number"`
		CreatedId   int       `db:"created_id"`
		CreatedAt   time.Time `db:"created_at"`
	}

	projectPersonModel struct {
		table string
	}
)

func newProjectPersonModel() *projectPersonModel {
	return &projectPersonModel{
		table: "project_person",
	}
}

func (m *projectPersonModel) Create(projectPerson *ProjectPerson) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?)", m.table, insertProjectPersonStr)
	_, err := db.Exec(sqlStr, projectPerson.ProjectId, projectPerson.UserId, projectPerson.Type, projectPerson.Name, projectPerson.PhoneNumber, projectPerson.CreatedId)
	if err != nil {
		slog.Error("insert project person err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectPersonModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete project person err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectPersonModel) Find(id int) (*ProjectPerson, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
	projectPerson := new(ProjectPerson)
	if err := db.Get(projectPerson, sqlStr, id); err != nil {
		slog.Error("find project Person err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return projectPerson, nil
}

func (m *projectPersonModel) SelectByProjectId(projectId int) ([]*ProjectPerson, error) {
	sqlStr := fmt.Sprintf("select * from %s where `project_id` = ? order by `type` asc", m.table)
	var projectPersons []*ProjectPerson
	if err := db.Select(&projectPersons, sqlStr, projectId); err != nil {
		slog.Error("select project person err ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
		return nil, err
	}
	return projectPersons, nil
}

func (m *projectPersonModel) SelectByUserId(userId int) ([]*ProjectPerson, error) {
	sqlStr := fmt.Sprintf("select * from %s where `user_id` = ? ", m.table)
	var projectPersons []*ProjectPerson
	if err := db.Select(&projectPersons, sqlStr, userId); err != nil {
		slog.Error("select project person err ", "sql", sqlStr, "user_id", userId, "err ", err.Error())
		return nil, err
	}
	return projectPersons, nil
}

func (m *projectPersonModel) BatchInsert(projectPersons []*ProjectPerson) error {
	if len(projectPersons) == 0 {
		return nil
	}
	ph := strings.TrimRight(strings.Repeat("(?,?,?,?,?,?),", len(projectPersons)), ",")
	var args []interface{}
	for _, projectPerson := range projectPersons {
		args = append(args, projectPerson.ProjectId, projectPerson.UserId, projectPerson.Type, projectPerson.Name, projectPerson.PhoneNumber, projectPerson.CreatedId)
	}
	querySql := `insert into ` + m.table + `(` + insertProjectPersonStr + `) values ` + ph
	_, err := db.Exec(querySql, args...)
	if err != nil {
		slog.Error("batch insert project person error", "err", err)
	}
	return err
}

func (m *projectPersonModel) FindByProjectIdAndUserId(projectId, userId int) (*ProjectPerson, error) {
	sqlStr := fmt.Sprintf("select * from %s where `project_id` = ? and `user_id` = ? limit 1", m.table)
	projectPerson := new(ProjectPerson)
	if err := db.Get(projectPerson, sqlStr, projectId, userId); err != nil {
		slog.Error("find project person err ", "sql", sqlStr, "project_id ", projectId, "user_id", userId, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return projectPerson, nil
}

func (m *projectPersonModel) FindFirst(projectId int) (*ProjectPerson, error) {
	sqlStr := fmt.Sprintf("select * from %s where `project_id` = ? and `type` = ? limit 1", m.table)
	projectPerson := new(ProjectPerson)
	if err := db.Get(projectPerson, sqlStr, projectId, ProjectPersonTypeFirst); err != nil {
		slog.Error("find project first person err ", "sql", sqlStr, "project_id ", projectId, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return projectPerson, nil
}

func (m *projectPersonModel) SelectByProjectIds(projectIds []int) ([]*ProjectPerson, error) {
	projectPersons := make([]*ProjectPerson, 0)
	if len(projectIds) == 0 {
		return projectPersons, nil
	}

	sqlStr := fmt.Sprintf("select * from %s where `type` = ? and id in (?) order by created_at desc", m.table)
	query, args, err := sqlx.In(sqlStr, ProjectPersonTypeFirst, projectIds)
	if err != nil {
		slog.Error("select  project person by projectIds error", "sql", sqlStr, "projectIds", projectIds, "err ", err.Error())
		return nil, err
	}

	if err := db.Select(&projectPersons, query, args...); err != nil {
		slog.Error("select  project person by projectIds error", "sql", sqlStr, "projectIds", projectIds, "err ", err.Error())
		return nil, err
	}
	return projectPersons, nil
}

func (m *projectPersonModel) SelectGroupCountByUserIds(userIds []int) ([]*UserCountItem, error) {
	items := make([]*UserCountItem, 0)
	if len(userIds) == 0 {
		return items, nil
	}

	sqlStr := fmt.Sprintf("select `user_id`,count(*) as `count` from %s where user_id in (?) group by user_id ", m.table)
	query, args, err := sqlx.In(sqlStr, userIds)
	if err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "err", err.Error())
		return nil, err
	}

	if err := db.Select(&items, query, args...); err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "err", err.Error())
		return nil, err
	}
	return items, nil
}

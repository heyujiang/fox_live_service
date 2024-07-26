package model

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"time"
)

const (
	inertProjectPersonStr = "`project_id`,`user_id`,`name`,`phone_number`,`type`,`created_id`,`updated_id`"
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
		UpdatedId   int       `db:"updated_id"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
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
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?)", m.table, inertProjectPersonStr)
	_, err := db.Exec(sqlStr, projectPerson.ProjectId, projectPerson.UserId, projectPerson.Name, projectPerson.PhoneNumber, projectPerson.Type, projectPerson.CreatedId, projectPerson.UpdatedId)
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

package model

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"strings"
	"time"
)

const (
	insertProjectContactStr = "`project_id`,`name`,`phone_number`,`type`,`description`,`created_id`"
)

const (
	ProjectContactTypeFirst  = iota + 1 // 第一负责人
	ProjectContactTypeSecond            // 第二负责人
	ProjectContactTypeCommon            // 项目成员
)

var (
	ProjectContactModel = newProjectContactModel()
)

type (
	ProjectContact struct {
		Id          int       `db:"id"`
		ProjectId   int       `db:"project_id"`
		Type        int       `db:"type"`
		Name        string    `db:"name"`
		PhoneNumber string    `db:"phone_number"`
		Description string    `db:"description"`
		CreatedId   int       `db:"created_id"`
		CreatedAt   time.Time `db:"created_at"`
	}

	projectContactModel struct {
		table string
	}
)

func newProjectContactModel() *projectContactModel {
	return &projectContactModel{
		table: "project_contact",
	}
}

func (m *projectContactModel) Create(projectContact *ProjectContact) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?)", m.table, insertProjectContactStr)
	_, err := db.Exec(sqlStr, projectContact.ProjectId, projectContact.Name, projectContact.PhoneNumber, projectContact.Type, projectContact.Description, projectContact.CreatedId)
	if err != nil {
		slog.Error("insert project contact err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectContactModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete project contact err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectContactModel) Find(id int) (*ProjectContact, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
	projectContact := new(ProjectContact)
	if err := db.Get(projectContact, sqlStr, id); err != nil {
		slog.Error("find project contact err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return projectContact, nil
}

func (m *projectContactModel) SelectByProjectId(projectId int) ([]*ProjectContact, error) {
	sqlStr := fmt.Sprintf("select * from %s where `project_id` = ? order by `type` asc", m.table)
	var projectContacts []*ProjectContact
	if err := db.Select(&projectContacts, sqlStr, projectId); err != nil {
		slog.Error("select project contact err ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
		return nil, err
	}
	return projectContacts, nil
}

func (m *projectContactModel) BatchInsert(projectContacts []*ProjectContact) error {
	if len(projectContacts) == 0 {
		return nil
	}
	ph := strings.TrimRight(strings.Repeat("(?,?,?,?,?,?),", len(projectContacts)), ",")
	var args []interface{}
	for _, projectContact := range projectContacts {
		args = append(args, projectContact.ProjectId, projectContact.Name, projectContact.PhoneNumber, projectContact.Type, projectContact.Description, projectContact.CreatedId)
	}
	querySql := `insert into ` + m.table + `(` + insertProjectContactStr + `) values ` + ph
	_, err := db.Exec(querySql, args...)
	if err != nil {
		slog.Error("batch insert project contact error", "err", err, "sql", querySql, "args", args)
	}
	return err
}

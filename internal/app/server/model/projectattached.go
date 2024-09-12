package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
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
		IsDeleted int       `db:"is_deleted"`
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

func (m *projectAttachedModel) Delete(id, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `is_deleted` = ? , `updated_id` = ? where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, ProjectDeletedYes, uid)
	if err != nil {
		slog.Error("delete project attached err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectAttachedModel) DeleteByProjectId(projectId int) error {
	sqlStr := fmt.Sprintf("update %s set `is_deleted` = ? where `project_id` = %d", m.table, projectId)
	_, err := db.Exec(sqlStr, ProjectDeletedYes)
	if err != nil {
		slog.Error("delete project attached by project id err ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectAttachedModel) Find(id int) (*ProjectAttached, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? and `is_deleted` = ? limit 1", m.table)
	projectAttached := new(ProjectAttached)
	args := []interface{}{id, ProjectDeletedNo}
	if err := db.Get(projectAttached, sqlStr, args...); err != nil {
		slog.Error("find project attached err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return projectAttached, nil
}

func (m *projectAttachedModel) GetProjectAttachedByCond(cond *ProjectAttachedCond) ([]*ProjectAttached, error) {
	sqlCond, args := m.buildProjectAttachedCond(cond)
	args = append([]interface{}{ProjectDeletedNo}, args...)
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? %s ", m.table, sqlCond)
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
		sqlCond += " and id = ? "
		args = append(args, cond.Id)
	}
	if cond.ProjectId > 0 {
		sqlCond += " and project_id = ? "
		args = append(args, cond.ProjectId)
	}
	if cond.RecordId > 0 {
		sqlCond += " and record_id = ? "
		args = append(args, cond.RecordId)
	}
	if cond.NodeId > 0 {
		sqlCond += " and node_id = ? "
		args = append(args, cond.NodeId)
	}
	if cond.UserId > 0 {
		sqlCond += " and user_id = ? "
		args = append(args, cond.UserId)
	}
	return
}

func (m *projectAttachedModel) GetAllByProjectId(projectId int) ([]*ProjectAttached, error) {
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and project_id = ?", m.table)
	var projectAttacheds []*ProjectAttached
	args := []interface{}{ProjectDeletedNo, projectId}
	if err := db.Select(&projectAttacheds, sqlStr, args...); err != nil {
		slog.Error("get project attached error ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
		return nil, err
	}
	return projectAttacheds, nil
}

func (m *projectAttachedModel) SelectGroupCountByUserIds(userIds []int, startTime, endTime time.Time) ([]*UserCountItem, error) {
	items := make([]*UserCountItem, 0)
	if len(userIds) == 0 {
		return items, nil
	}

	sqlStr := fmt.Sprintf("select `user_id`,count(*) as `count` from %s where `is_deleted` = ? and `created_at` >= ? and `created_at` < ? and user_id in (?) group by user_id ", m.table)
	query, args, err := sqlx.In(sqlStr, ProjectDeletedNo, startTime, endTime, userIds)
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

// GetAllByProjectIdsAndUserId 根据项目id数组（projectIds） 和 创建时间范围（createdAts）查询 指定用户（userId）的提交记录
func (m *projectAttachedModel) GetAllByProjectIdsAndUserId(projectIds []int, userId int, createdAts []*time.Time) ([]*ProjectAttached, error) {
	var projectAttached []*ProjectAttached

	if len(projectIds) == 0 {
		return projectAttached, nil
	}
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and project_id in (?) and `user_id` = ? and `created_at` >= ? and `created_at` < ? order by created_at desc", m.table)
	query, args, err := sqlx.In(sqlStr, ProjectDeletedNo, projectIds, userId, createdAts[0], createdAts[1])
	if err != nil {
		slog.Error("get all project attached by project_ids and user_id error ", "sql", sqlStr, "projectIds", projectIds, "userId", userId, "createAts", createdAts, "err", err.Error())
		return nil, err
	}

	if err := db.Select(&projectAttached, query, args...); err != nil {
		slog.Error("get all project attached by project_ids and user_id error ", "sql", sqlStr, "projectIds", projectIds, "userId", userId, "createAts", createdAts, "err", err.Error())
		return nil, err
	}
	return projectAttached, nil
}

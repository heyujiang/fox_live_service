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
	ProjectRecordStateFinished int = iota + 1
	ProjectRecordStateIng

	latestRecordCountView = 5

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
		IsDeleted   int       `db:"is_deleted"`
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
	res, err := db.Exec(sqlStr, projectRecord.ProjectId, projectRecord.ProjectName, projectRecord.NodeId, projectRecord.NodeName, projectRecord.UserId, projectRecord.Username, projectRecord.Overview, projectRecord.State, projectRecord.CreatedId, projectRecord.UpdatedId)
	if err != nil {
		slog.Error("insert project record err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastInsertId, _ := res.LastInsertId()
	projectRecord.Id = cast.ToInt(lastInsertId)
	return nil
}

func (m *projectRecordModel) Delete(id, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `is_deleted` = ? , `updated_id` = ? where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, ProjectDeletedYes, uid)
	if err != nil {
		slog.Error("delete project record err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectRecordModel) DeleteByProjectId(projectId, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `is_deleted` = ? , `updated_id` = ? where `project_id` = %d", m.table, projectId)
	_, err := db.Exec(sqlStr, ProjectDeletedYes, uid)
	if err != nil {
		slog.Error("delete project record by project id err ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
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
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? and `is_deleted` = ? limit 1", m.table)
	projectRecord := new(ProjectRecord)
	if err := db.Get(projectRecord, sqlStr, id, ProjectDeletedNo); err != nil {
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
	sqlStr := fmt.Sprintf("select count(*) from %s where `is_deleted` = ? %s order by created_at desc", m.table, sqlCond)
	var count int
	args = append([]interface{}{ProjectDeletedNo}, args...)
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
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? %s order by created_at desc limit %d,%d", m.table, sqlCond, (pageIndex-1)*pageSize, pageSize)
	var projectRecords []*ProjectRecord
	args = append([]interface{}{ProjectDeletedNo}, args...)
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
		sqlCond += " and id = ? "
		args = append(args, cond.Id)
	}
	if cond.ProjectId > 0 {
		sqlCond += " and project_id = ? "
		args = append(args, cond.ProjectId)
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

func (m *projectRecordModel) GetAllByProjectId(projectId int) ([]*ProjectRecord, error) {
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and project_id = ?", m.table)
	var projectRecords []*ProjectRecord
	if err := db.Select(&projectRecords, sqlStr, ProjectDeletedNo, projectId); err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
		return nil, err
	}
	return projectRecords, nil
}

func (m *projectRecordModel) SelectByUserId(userId int) ([]*ProjectRecord, error) {
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and user_id = ? order by created_at desc limit %d", m.table, latestRecordCountView)
	var projectRecords []*ProjectRecord
	if err := db.Select(&projectRecords, sqlStr, ProjectDeletedNo, userId); err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "user_id", userId, "err ", err.Error())
		return nil, err
	}
	return projectRecords, nil
}

func (m *projectRecordModel) Select() ([]*ProjectRecord, error) {
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? order by created_at desc limit %d", m.table, latestRecordCountView*2)
	var projectRecords []*ProjectRecord
	if err := db.Select(&projectRecords, sqlStr, ProjectDeletedNo); err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return projectRecords, nil
}

func (m *projectRecordModel) SelectProjectIdFromCreatedAt(at *time.Time) ([]int, error) {
	sqlStr := fmt.Sprintf("select `project_id` from %s where `is_deleted` = ? and created_at > ? group by `project_id` ", m.table)
	var projectIds []int
	if err := db.Select(&projectIds, sqlStr, ProjectDeletedNo, at); err != nil {
		slog.Error("get project ids from created at error ", "sql", sqlStr, "at", at, "err ", err.Error())
		return nil, err
	}
	return projectIds, nil
}

func (m *projectRecordModel) SelectGroupCountByUserIds(userIds []int) ([]*UserCountItem, error) {
	items := make([]*UserCountItem, 0)
	if len(userIds) == 0 {
		return items, nil
	}

	sqlStr := fmt.Sprintf("select `user_id`,count(*) as `count` from %s where `is_deleted` = ? and user_id in (?) group by user_id ", m.table)
	query, args, err := sqlx.In(sqlStr, ProjectDeletedNo, userIds)
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

func (m *projectRecordModel) SelectByUserIds(userIds []int) ([]*ProjectRecord, error) {
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and user_id in (?) order by created_at desc limit %d", m.table, latestRecordCountView)
	query, args, err := sqlx.In(sqlStr, ProjectDeletedNo, userIds)
	if err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "err", err.Error())
		return nil, err
	}

	var projectRecords []*ProjectRecord
	if err := db.Select(&projectRecords, query, args...); err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "userIds", userIds, "err ", err.Error())
		return nil, err
	}
	return projectRecords, nil
}

func (m *projectRecordModel) SelectByProjectIds(projectIds []int) ([]*ProjectRecord, error) {
	var projectRecords []*ProjectRecord

	if len(projectIds) == 0 {
		return projectRecords, nil
	}
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and project_id in (?) order by created_at desc limit %d", m.table, latestRecordCountView)
	query, args, err := sqlx.In(sqlStr, ProjectDeletedNo, projectIds)
	if err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "err", err.Error())
		return nil, err
	}

	if err := db.Select(&projectRecords, query, args...); err != nil {
		slog.Error("get project record error ", "sql", sqlStr, "userIds", projectIds, "err ", err.Error())
		return nil, err
	}
	return projectRecords, nil
}

// GetAllByProjectIdsAndUserId 根据项目id数组（projectIds） 和 创建时间范围（createdAts）查询 指定用户（userId）的提交记录
func (m *projectRecordModel) GetAllByProjectIdsAndUserId(projectIds []int, userId int, createdAts []*time.Time) ([]*ProjectRecord, error) {
	var projectRecords []*ProjectRecord

	if len(projectIds) == 0 {
		return projectRecords, nil
	}
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and project_id in (?) and `user_id` = ? and `created_at` >= ? and `created_at` < ? order by created_at desc", m.table)
	query, args, err := sqlx.In(sqlStr, ProjectDeletedNo, projectIds, userId, createdAts[0], createdAts[1])
	if err != nil {
		slog.Error("get all project record by project_ids and user_id error ", "sql", sqlStr, "projectIds", projectIds, "userId", userId, "createAts", createdAts, "err", err.Error())
		return nil, err
	}

	if err := db.Select(&projectRecords, query, args...); err != nil {
		slog.Error("get all project record by project_ids and user_id error ", "sql", sqlStr, "projectIds", projectIds, "userId", userId, "createAts", createdAts, "err", err.Error())
		return nil, err
	}
	return projectRecords, nil
}

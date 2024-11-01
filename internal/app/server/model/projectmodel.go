package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
	"strings"
	"time"
)

const (
	ProjectAttrCentralized   = iota + 1 // 集中式
	ProjectAttrDistributed              // 分布式
	ProjectAttrDecentralized            // 分散式
)

const (
	ProjectStateWait      = iota + 1 // 待定
	ProjectStateRecommend            // 推荐
	ProjectStateFinished             // 已完成
	ProjectStateStop                 // 终止
)

const (
	ProjectWaitAudit = iota
	ProjectAudited
)

const (
	ProjectTypeWind                   = iota + 1 // 风电
	ProjectTypeLight                             // 光伏
	ProjectTypeWindAndLight                      // 风电+光伏
	ProjectTypeStorage                           // 储能
	ProjectTypeWindAndStorage                    // 风电+储能
	ProjectTypeLightAndStorage                   // 光伏+储能
	ProjectTypeWindAndLightAndStorage            // 风光储一体
)

const (
	ProjectDeletedNo = iota
	ProjectDeletedYes
)

const (
	myProjectViewCountView = 8
	inertProjectStr        = "`name`,`description`,`attr`,`state`,`type`,`node_id`,`node_name`,`schedule`,`capacity`,`properties`," +
		"`area`,`address`,`connect`,`investment_agreement`,`business_condition`,`star`,`user_id`,`username`,`begin_time`,`created_id`,`updated_id`"
)

var (
	ProjectModel = newProjectModel()

	ProjectAttrDesc = map[int]string{
		ProjectAttrCentralized:   "集中式",
		ProjectAttrDistributed:   "分布式",
		ProjectAttrDecentralized: "分散式",
	}

	ProjectStateDesc = map[int]string{
		ProjectStateWait:      "待定",
		ProjectStateRecommend: "推荐",
		ProjectStateFinished:  "已完成",
		ProjectStateStop:      "终止",
	}

	ProjectTypeDesc = map[int]string{
		ProjectTypeWind:                   "风电",
		ProjectTypeLight:                  "光伏",
		ProjectTypeWindAndLight:           "风电+光伏",
		ProjectTypeStorage:                "储能",
		ProjectTypeWindAndStorage:         "风电+储能",
		ProjectTypeLightAndStorage:        "光伏+储能",
		ProjectTypeWindAndLightAndStorage: "风光储一体",
	}
)

type (
	Project struct {
		Id                  int       `db:"id"`
		Name                string    `db:"name"`
		Description         string    `db:"description"`
		Attr                int       `db:"attr"`
		State               int       `db:"state"`
		Type                int       `db:"type"`
		NodeId              int       `db:"node_id"`
		NodeName            string    `db:"node_name"`
		Schedule            float64   `db:"schedule"`
		Capacity            float64   `db:"capacity"`
		Properties          string    `db:"properties"`
		Area                float64   `db:"area"`
		Address             string    `db:"address"`
		Connect             string    `db:"connect"`
		Star                int       `db:"star"`
		UserId              int       `db:"user_id"`
		Username            string    `db:"username"`
		InvestmentAgreement string    `db:"investment_agreement"`
		BusinessCondition   string    `db:"business_condition"`
		BeginTime           time.Time `db:"begin_time"`
		IsAudit             int       `db:"is_audit"`
		IsDeleted           int       `db:"is_deleted"`
		CreatedId           int       `db:"created_id"`
		UpdatedId           int       `db:"updated_id"`
		CreatedAt           time.Time `db:"created_at"`
		UpdatedAt           time.Time `db:"updated_at"`
	}

	projectModel struct {
		table string
	}

	ProjectCond struct {
		Name      string
		UserId    int
		CreatedAt []time.Time
		Star      int
		Type      int
		Ids       []int
		State     int
		IsAudit   int
	}
)

func newProjectModel() *projectModel {
	return &projectModel{
		table: "project",
	}
}

func (m *projectModel) Create(project *Project) (int, error) {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", m.table, inertProjectStr)
	res, err := db.Exec(sqlStr, project.Name, project.Description, project.Attr, project.State, project.Type, project.NodeId,
		project.NodeName, project.Schedule, project.Capacity, project.Properties, project.Area, project.Address, project.Connect,
		project.InvestmentAgreement, project.BusinessCondition, project.Star, project.UserId, project.Username, project.BeginTime,
		project.CreatedId, project.UpdatedId)
	if err != nil {
		slog.Error("insert project err ", "sql", sqlStr, "err ", err.Error())
		return 0, err
	}
	projectId, _ := res.LastInsertId()
	return cast.ToInt(projectId), nil
}

func (m *projectModel) Delete(id int, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `is_deleted` = ? , `updated_id` = ? where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, ProjectDeletedYes, uid)
	if err != nil {
		slog.Error("delete project err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectModel) Update(project *Project) error {
	sqlStr := fmt.Sprintf("update %s set `name` = ? , `description` = ?, `attr`= ? , `star` = ?, `state` = ? , `type` = ?, `capacity` = ? , `properties`= ?, `area` = ?, `address` = ?, `connect` = ?, `investment_agreement` = ?, `business_condition` =? , `begin_time` = ?, `updated_id`= ? where `id` = %d", m.table, project.Id)
	_, err := db.Exec(sqlStr, project.Name, project.Description, project.Attr, project.Star, project.State, project.Type, project.Capacity, project.Properties, project.Area, project.Address, project.Connect, project.InvestmentAgreement, project.BusinessCondition, project.BeginTime, project.UpdatedId)
	if err != nil {
		slog.Error("update project err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectModel) UpdateAfterRecord(project *Project) error {
	sqlStr := fmt.Sprintf("update %s set `node_id` = ? , `node_name` = ?, `schedule`= ?, `updated_id`= ? where `id` = %d", m.table, project.Id)
	_, err := db.Exec(sqlStr, project.NodeId, project.NodeName, project.Schedule, project.UpdatedId)
	if err != nil {
		slog.Error("update project after create record err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectModel) UpdateFirstPerson(projectId int, userId int, userName string, uid int) error {
	sqlStr := fmt.Sprintf("update %s set  `user_id` = ? , `username` = ?, `updated_id`= ? where `id` = %d", m.table, projectId)
	_, err := db.Exec(sqlStr, userId, userName, uid)
	if err != nil {
		slog.Error("update project err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectModel) UpdateSchedule(projectId int, schedule float64, uid int) error {
	sqlStr := fmt.Sprintf("update %s set  `schedule` = ? , `updated_id`= ? where `id` = %d", m.table, projectId)
	_, err := db.Exec(sqlStr, schedule, uid)
	if err != nil {
		slog.Error("update project schedule err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectModel) Find(id int) (*Project, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? and `is_deleted` = ? limit 1", m.table)
	project := new(Project)
	if err := db.Get(project, sqlStr, id, ProjectDeletedNo); err != nil {
		slog.Error("find project err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return project, nil
}

func (m *projectModel) Select() ([]*Project, error) {
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? ", m.table)
	projects := make([]*Project, 0)
	if err := db.Select(&projects, sqlStr, ProjectDeletedNo); err != nil {
		slog.Error("select project err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return projects, nil
}

func (m *projectModel) GetProjectByCond(cond *ProjectCond, sortCond string, pageIndex, pageSize int) ([]*Project, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	sqlCond, args := m.buildProjectCond(cond)
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ?  %s order by %s created_at desc limit %d,%d", m.table, sqlCond, sortCond, (pageIndex-1)*pageSize, pageSize)

	var projects []*Project
	if err := db.Select(&projects, sqlStr, append([]interface{}{ProjectDeletedNo}, args...)...); err != nil {
		slog.Error("get projects error ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return projects, nil
}

func (m *projectModel) GetAllProjectByCond(cond *ProjectCond) ([]*Project, error) {
	sqlCond, args := m.buildProjectCond(cond)
	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ?  %s order by created_at desc", m.table, sqlCond)

	var projects []*Project
	if err := db.Select(&projects, sqlStr, append([]interface{}{ProjectDeletedNo}, args...)...); err != nil {
		slog.Error("get projects error ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return projects, nil
}

func (m *projectModel) GetProjectCountByCond(cond *ProjectCond) (int, error) {
	sqlCond, args := m.buildProjectCond(cond)
	sqlStr := fmt.Sprintf("select count(*) from %s where `is_deleted` = ? %s order by created_at desc", m.table, sqlCond)
	var count int
	if err := db.Get(&count, sqlStr, append([]interface{}{ProjectDeletedNo}, args...)...); err != nil {
		slog.Error("get project count error ", "sql", sqlStr, "err ", err.Error())
		return 0, err
	}
	return count, nil
}

func (m *projectModel) buildProjectCond(cond *ProjectCond) (sqlCond string, args []interface{}) {
	if cond == nil {
		return
	}

	if len(cond.Ids) > 0 {
		sqlCond += fmt.Sprintf(" and id in (%s) ", strings.Repeat(",?", len(cond.Ids))[1:])
		argIds := make([]interface{}, 0, len(cond.Ids))
		for _, id := range cond.Ids {
			argIds = append(argIds, id)
		}
		args = append(args, argIds...)
	}

	if cond.UserId > 0 {
		sqlCond += " and user_id = ?"
		args = append(args, cond.UserId)
	}

	if cond.Name != "" {
		sqlCond += " and name like ?"
		args = append(args, "%"+cond.Name+"%")
	}

	if cond.Star > 0 {
		sqlCond += " and star = ?"
		args = append(args, cond.Star)
	}

	if cond.Type > 0 {
		sqlCond += " and type = ?"
		args = append(args, cond.Type)
	}

	if len(cond.CreatedAt) == 2 {
		sqlCond += " and created_at >= ? and created_at <= ?"
		args = append(args, cond.CreatedAt[0], cond.CreatedAt[1])
	}

	if cond.State != 0 {
		sqlCond += " and state = ?"
		args = append(args, cond.State)
	}

	if cond.IsAudit == ProjectAudited {
		sqlCond += " and is_audit = ?"
		args = append(args, ProjectAudited)
	}

	return
}

func (m *projectModel) UpdateProjectState(id int, state, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `state` = ? , `updated_id`= ?  where `id` = %d and `is_deleted` = %d", m.table, id, ProjectDeletedNo)
	_, err := db.Exec(sqlStr, state, uid)
	if err != nil {
		slog.Error("update project state err ", "sql", sqlStr, "state", state, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectModel) SelectByIds(ids []int) ([]*Project, error) {
	projects := make([]*Project, 0)
	if len(ids) == 0 {
		return projects, nil
	}

	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and id in (?) order by created_at ", m.table)
	query, args, err := sqlx.In(sqlStr, ProjectDeletedNo, ids)

	if err != nil {
		slog.Error("select my project by ids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}

	if err := db.Select(&projects, query, args...); err != nil {
		slog.Error("select my project by ids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}
	return projects, nil
}

func (m *projectModel) SelectLatestProject(beginTime *time.Time, notExistIds, existIds []int, star int) ([]*Project, error) {
	args := make([]interface{}, 0)
	sqlStr := " and created_at > ? "
	args = append(args, ProjectDeletedNo, beginTime)

	if star == 3 {
		sqlStr += " and star = ? "
		args = append(args, star)
	}

	if len(notExistIds) > 0 {
		notInStr := strings.TrimRight(strings.Repeat("?,", len(notExistIds)), ",")
		sqlStr += " and id not in ( " + notInStr + " ) "

		notIds := make([]interface{}, 0)
		for _, id := range notExistIds {
			notIds = append(notIds, id)
		}
		args = append(args, notIds...)
	}

	if len(existIds) > 0 {
		notInStr := strings.TrimRight(strings.Repeat("?,", len(existIds)), ",")
		sqlStr += " and id in ( " + notInStr + " ) "

		notIds := make([]interface{}, 0)
		for _, id := range existIds {
			notIds = append(notIds, id)
		}
		args = append(args, notIds...)
	}

	query := fmt.Sprintf("select * from %s where `is_deleted` = ? %s order by `created_at` desc", m.table, sqlStr)

	projects := make([]*Project, 0)
	if err := db.Select(&projects, query, args...); err != nil {
		slog.Error("select my project by ids error", "sql", query, "sqlStr", sqlStr, "args", args, "err ", err.Error())
		return projects, err
	}
	return projects, nil
}

func (m *projectModel) GetAllNoFinished() ([]*Project, error) {
	projects := make([]*Project, 0)

	sqlStr := fmt.Sprintf("select * from %s where `is_deleted` = ? and `state` in (? , ?) order by created_at ", m.table)
	query, args, err := sqlx.In(sqlStr, ProjectDeletedNo, ProjectStateWait, ProjectStateRecommend)

	if err != nil {
		slog.Error("get all not finished project error", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}

	if err := db.Select(&projects, query, args...); err != nil {
		slog.Error("get all not finished project error", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return projects, nil
}

func (m *projectModel) Audit(id int, isAudit int, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `is_audit` = ? , `updated_id`= ?  where `id` = %d and `is_deleted` = %d", m.table, id, ProjectDeletedNo)
	_, err := db.Exec(sqlStr, isAudit, uid)
	if err != nil {
		slog.Error("audit project err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

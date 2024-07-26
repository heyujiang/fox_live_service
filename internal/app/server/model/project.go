package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
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
	ProjectStateStop                 // 终止
	ProjectStateFinished             // 已完成
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
	inertProjectStr = "`name`,`description`,`attr`,`state`,`type`,`node_id`,`node_name`,`schedule`,`capacity`,`properties`," +
		"`area`,`address`,`connect`,`investment_agreement`,`business_condition`,`begin_time`,`created_id`,`updated_id`"
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
		ProjectStateStop:      "终止",
		ProjectStateFinished:  "已完成",
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
		InvestmentAgreement string    `db:"investment_agreement"`
		BusinessCondition   string    `db:"business_condition"`
		BeginTime           time.Time `db:"begin_time"`
		CreatedId           int       `db:"created_id"`
		UpdatedId           int       `db:"updated_id"`
		CreatedAt           time.Time `db:"created_at"`
		UpdatedAt           time.Time `db:"updated_at"`
	}

	projectModel struct {
		table string
	}

	ProjectCond struct {
		Id     int
		Name   string
		NodeId int
		Attr   int
		State  int
		Type   int
	}
)

func newProjectModel() *projectModel {
	return &projectModel{
		table: "project",
	}
}

func (m *projectModel) Create(project *Project) (int, error) {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", m.table, inertProjectStr)
	res, err := db.Exec(sqlStr, project.Name, project.Description, project.Attr, project.State, project.Type, project.NodeId,
		project.NodeName, project.Schedule, project.Capacity, project.Properties, project.Area, project.Address, project.Connect,
		project.InvestmentAgreement, project.BusinessCondition, project.BeginTime, project.CreatedId, project.UpdatedId)
	if err != nil {
		slog.Error("insert project err ", "sql", sqlStr, "err ", err.Error())
		return 0, err
	}
	projectId, _ := res.LastInsertId()
	return cast.ToInt(projectId), nil
}

func (m *projectModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete project err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectModel) Update(project *Project) error {
	sqlStr := fmt.Sprintf("update %s set `name` = ? ,  `attr`= ? , `type` = ?, `updated_id`= ? where `id` = %d", m.table, project.Id)
	_, err := db.Exec(sqlStr, project.Name, project.Attr, project.Type, project.UpdatedId)
	if err != nil {
		slog.Error("update project err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectModel) Find(id int) (*Project, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
	project := new(Project)
	if err := db.Get(project, sqlStr, id); err != nil {
		slog.Error("find project err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return project, nil
}

func (m *projectModel) GetProjectByCond(cond *ProjectCond, pageIndex, pageSize int) ([]*Project, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	sqlCond, args := m.buildProjectCond(cond)
	sqlStr := fmt.Sprintf("select * from %s where 1 = 1 %s limit %d,%d", m.table, sqlCond, (pageIndex-1)*pageSize, pageSize)
	var projects []*Project
	if err := db.Select(&projects, sqlStr, args...); err != nil {
		slog.Error("get projects error ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return projects, nil
}

func (m *projectModel) GetProjectCountByCond(cond *ProjectCond) (int, error) {
	sqlCond, args := m.buildProjectCond(cond)
	sqlStr := fmt.Sprintf("select count(*) from %s where 1 = 1 %s", m.table, sqlCond)
	var count int
	if err := db.Get(&count, sqlStr, args...); err != nil {
		slog.Error("get project count error ", "sql", sqlStr, "err ", err.Error())
		return 0, err
	}
	return count, nil
}

func (m *projectModel) buildProjectCond(cond *ProjectCond) (sqlCond string, args []interface{}) {
	if cond == nil {
		return
	}

	if cond.Id > 0 {
		sqlCond += "and id = ?"
		args = append(args, cond.Id)
	}

	if cond.Name != "" {
		sqlCond += " and name = ?"
		args = append(args, cond.Name)
	}

	if cond.NodeId > 0 {
		sqlCond += " and node_id = ?"
		args = append(args, cond.Name)
	}

	if _, ok := ProjectAttrDesc[cond.Attr]; ok {
		sqlCond += " and attr = ?"
		args = append(args, cond.Attr)
	}

	if _, ok := ProjectStateDesc[cond.State]; ok {
		sqlCond += " and state = ?"
		args = append(args, cond.State)
	}

	if _, ok := ProjectTypeDesc[cond.Type]; ok {
		sqlCond += " and type = ?"
		args = append(args, cond.Type)
	}

	return
}

func (m *projectModel) UpdateProjectState(id int, state, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `state` = ? , `updated_id`= ?  where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, state, uid)
	if err != nil {
		slog.Error("update project state err ", "sql", sqlStr, "state", state, "err ", err.Error())
		return err
	}
	return nil
}

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
	insertProjectNodeStr = "`project_id`,`p_id`,`node_id`,`name`,`is_leaf`,`sort`,`state`,`created_id`,`updated_id`"
)

const (
	ProjectNodeStateWaitBegin = iota + 1 // 未开始
	ProjectNodeStateInProcess            // 进行中
	ProjectNodeStateFinished             // 已完成
)

var (
	ProjectNodeModel = newProjectNodeModel()

	ProjectNodeStateDesc = map[int]string{
		ProjectNodeStateWaitBegin: "未开始",
		ProjectNodeStateInProcess: "进行中",
		ProjectNodeStateFinished:  "已完成",
	}
)

type (
	ProjectNode struct {
		Id        int       `db:"id"`
		ProjectId int       `db:"project_id"`
		PId       int       `db:"p_id"` // node pid
		NodeId    int       `db:"node_id"`
		Name      string    `db:"name"` // 节点名称
		IsLeaf    int       `db:"is_leaf"`
		Sort      int       `db:"sort"`
		State     int       `db:"state"`
		CreatedId int       `db:"created_id"`
		UpdatedId int       `db:"updated_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	projectNodeModel struct {
		table string
	}

	ProjectNodeCond struct {
		Id        int
		ProjectId int
		NodeId    int
		UserId    int
	}
)

func newProjectNodeModel() *projectNodeModel {
	return &projectNodeModel{
		table: "project_node",
	}
}

func (m *projectNodeModel) Create(projectNode *ProjectNode) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?)", m.table, insertProjectNodeStr)
	_, err := db.Exec(sqlStr, projectNode.ProjectId, projectNode.PId, projectNode.NodeId, projectNode.Name, projectNode.IsLeaf, projectNode.Sort, projectNode.State, projectNode.CreatedId, projectNode.UpdatedId)
	if err != nil {
		slog.Error("insert project Node err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectNodeModel) Find(id int) (*ProjectNode, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1 ", m.table)
	projectNode := new(ProjectNode)
	if err := db.Get(projectNode, sqlStr, id); err != nil {
		slog.Error("find project Node err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return projectNode, nil
}

func (m *projectNodeModel) UpdateProjectNodeState(id int, state, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `state` = ? , `updated_id`= ?  where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, state, uid)
	if err != nil {
		slog.Error("update project node state err ", "sql", sqlStr, "state", state, "err ", err.Error())
		return err
	}
	return nil
}

func (m *projectNodeModel) BatchInsert(projectNodes []*ProjectNode) error {
	if len(projectNodes) == 0 {
		return nil
	}
	ph := strings.TrimRight(strings.Repeat("(?,?,?,?,?,?,?,?,?),", len(projectNodes)), ",")
	var args []interface{}
	for _, projectNode := range projectNodes {
		args = append(args, projectNode.ProjectId, projectNode.PId, projectNode.NodeId, projectNode.Name, projectNode.IsLeaf, projectNode.Sort, projectNode.State, projectNode.CreatedId, projectNode.UpdatedId)
	}
	querySql := `insert into ` + m.table + `(` + insertProjectNodeStr + `) values ` + ph
	_, err := db.Exec(querySql, args...)
	if err != nil {
		slog.Error("batch insert project node error", "err", err)
	}
	return err
}

func (m *projectNodeModel) SelectByProjectId(projectId int) ([]*ProjectNode, error) {
	sqlStr := fmt.Sprintf("select * from %s where `project_id` = ? order by sort asc", m.table)
	projectNodes := make([]*ProjectNode, 0)
	if err := db.Select(&projectNodes, sqlStr, projectId); err != nil {
		slog.Error("select project node by project id err ", "sql", sqlStr, "project_id", projectId, "err ", err.Error())
		return nil, err
	}
	return projectNodes, nil
}

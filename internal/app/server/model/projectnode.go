package model

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"time"
)

const (
	inertProjectNodeStr = "`project_id`,`p_node_id`,`node_id`,`node_name`,`is_leaf`,`sort`,`state`,`created_at`,`updated_at`"
)

const (
	ProjectNodeStateWaitBegin = iota + 1 // 未开始
	ProjectNodeStateInProcess            // 未开始
	ProjectNodeStateFinished             // 未开始
)

var (
	ProjectNodeModel = newProjectNodeModel()

	ProjectNodeStateDesc = map[int]string{
		ProjectNodeStateWaitBegin: "未开始",
		ProjectNodeStateInProcess: "未开始",
		ProjectNodeStateFinished:  "未开始",
	}
)

type (
	ProjectNode struct {
		Id        int       `db:"id"`
		ProjectId string    `db:"project_id"`
		PNodeId   int       `json:"p_node_id"`
		NodeId    int       `db:"node_id"`
		NodeName  string    `db:"node_name"`
		IsLeaf    bool      `db:"is_leaf"`
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
		table: "project_Node",
	}
}

func (m *projectNodeModel) Create(projectNode *ProjectNode) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?)", m.table, inertProjectNodeStr)
	_, err := db.Exec(sqlStr, projectNode.ProjectId, projectNode.PNodeId, projectNode.NodeId, projectNode.NodeName, projectNode.IsLeaf, projectNode.Sort, projectNode.State, projectNode.CreatedId, projectNode.UpdatedId)
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

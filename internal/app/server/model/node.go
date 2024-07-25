package model

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"time"
)

const (
	NodeLeafYes = iota + 1
	NodeLeafNo

	inertNodeStr = "`name`,`pid`,`is_leaf`,`created_id`,`updated_id`"
)

var (
	NodeModel = newNodeModel()

	NodeLeafDesc = map[int]string{
		NodeLeafYes: "是",
		NodeLeafNo:  "否",
	}
)

type (
	Node struct {
		Id        int       `db:"id"`
		Name      string    `db:"name"`
		Pid       int       `db:"pid"`
		IsLeaf    int       `db:"is_leaf"`
		CreatedId int       `db:"created_id"`
		UpdatedId int       `db:"updated_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	nodeModel struct {
		table string
	}
)

func newNodeModel() *nodeModel {
	return &nodeModel{
		table: "node",
	}
}

// Insert 插入数据
func (m *nodeModel) Insert(node *Node) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?)", m.table, inertNodeStr)
	result, err := db.Exec(sqlStr, node.Name, node.Pid, node.IsLeaf, node.CreatedId, node.UpdatedId)
	if err != nil {
		slog.Error("insert node err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	node.Id = int(lastInsertId)
	return nil
}

// Delete 更新数据
func (m *nodeModel) Delete(id int) error {
	var err error
	tx, err := db.Begin()
	if err != nil {
		slog.Error("delete node error , begin tx error ", "err ", err.Error())
		if tx != nil {
			_ = tx.Rollback()
		}
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
		_ = tx.Commit()
	}()

	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err = db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete node error ", "sql", sqlStr, "id", id, "err ", err.Error())
		return err
	}

	sqlStr = fmt.Sprintf("delete from %s where `pid` = ? ", m.table)
	_, err = db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete node child error ", "sql", sqlStr, "pid", id, "err ", err.Error())
		return err
	}

	return nil
}

// Update 更新数据
func (m *nodeModel) Update(node *Node) error {
	sqlStr := fmt.Sprintf("update %s set `name` = ? ,`updated_id` = ? where `id` = %d", m.table, node.Id)
	_, err := db.Exec(sqlStr, node.Name, node.UpdatedId)
	if err != nil {
		slog.Error("update node err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

// Find 根据主键id单条查询
func (m *nodeModel) Find(id int) (*Node, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ?", m.table)
	node := new(Node)
	if err := db.Get(node, sqlStr, id); err != nil {
		slog.Error("find node err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return node, nil
}

// Select 查询所有数据
func (m *nodeModel) Select() ([]*Node, error) {
	sqlStr := fmt.Sprintf("select * from %s ", m.table)
	var nodes []*Node
	if err := db.Select(&nodes, sqlStr, []interface{}{}...); err != nil {
		slog.Error("select node err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return nodes, nil
}

// SelectByPid 根据PID查询所有子节点
func (m *nodeModel) SelectByPid(pid int) ([]*Node, error) {
	sqlStr := fmt.Sprintf("select * from %s where `pid` = ? ", m.table)
	var nodes []*Node
	if err := db.Select(&nodes, sqlStr, pid); err != nil {
		slog.Error("select node by pid error ", "sql", sqlStr, "pid", pid, "err ", err.Error())
		return nil, err
	}
	return nodes, nil
}

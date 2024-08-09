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

	inertNodeStr = "`name`,`pid`,`is_leaf`,`sort`,`created_id`,`updated_id`"
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
		Sort      int       `db:"sort"`
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
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?)", m.table, inertNodeStr)
	result, err := db.Exec(sqlStr, node.Name, node.Pid, node.IsLeaf, node.Sort, node.CreatedId, node.UpdatedId)
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
	tx := db.MustBegin()

	defer func() {
		if p := recover(); p != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("recover from %#v, rollback failed: %s", p, e)
			} else {
				err = fmt.Errorf("recover from %#v", p)
			}
		} else if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("transaction failed: %s, rollback failed: %s", err, e)
			}
		} else {
			err = tx.Commit()
		}
	}()

	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err = tx.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete node error ", "sql", sqlStr, "id", id, "err ", err.Error())
		return err
	}

	sqlStr = fmt.Sprintf("delete from %s where `pid` = ? ", m.table)
	_, err = tx.Exec(sqlStr, id)
	if err != nil {
		slog.Error("delete node child error ", "sql", sqlStr, "pid", id, "err ", err.Error())
		return err
	}

	return nil
}

// Update 更新数据
func (m *nodeModel) Update(node *Node) error {
	sqlStr := fmt.Sprintf("update %s set `name` = ? ,`sort` = ?,`updated_id` = ? where `id` = %d", m.table, node.Id)
	_, err := db.Exec(sqlStr, node.Name, node.Sort, node.UpdatedId)
	if err != nil {
		slog.Error("update node err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

// Find 根据主键id单条查询
func (m *nodeModel) Find(id int) (*Node, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
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
	sqlStr := fmt.Sprintf("select * from %s order by sort asc ,id asc ", m.table)
	var nodes []*Node
	if err := db.Select(&nodes, sqlStr, []interface{}{}...); err != nil {
		slog.Error("select node err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return nodes, nil
}

// SelectNotLeaf 查询所有非叶子结点
func (m *nodeModel) SelectNotLeaf() ([]*Node, error) {
	sqlStr := fmt.Sprintf("select * from %s where is_leaf = ? order by sort asc ,id asc ", m.table)
	var nodes []*Node
	if err := db.Select(&nodes, sqlStr, NodeLeafNo); err != nil {
		slog.Error("select leaf node err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return nodes, nil
}

// SelectByPid 根据PID查询所有子节点
func (m *nodeModel) SelectByPid(pid int) ([]*Node, error) {
	sqlStr := fmt.Sprintf("select * from %s where `pid` = ? order by sort asc ,id asc", m.table)
	var nodes []*Node
	if err := db.Select(&nodes, sqlStr, pid); err != nil {
		slog.Error("select node by pid error ", "sql", sqlStr, "pid", pid, "err ", err.Error())
		return nil, err
	}
	return nodes, nil
}

package model

import (
	"fmt"
	"golang.org/x/exp/slog"
)

const (
	insertFileStr = "`type`,`url`,`mime`,`size`,`filename`,`created_id`"
)

var (
	FileModel = newFileModel()
)

func newFileModel() *fileModel {
	return &fileModel{
		table: "file",
	}
}

type (
	File struct {
		Id        int    `db:"id"`
		Type      string `db:"type"`
		Url       string `db:"url"`
		Mime      string `db:"mime"`
		Size      int64  `db:"size"`
		Filename  string `db:"filename"`
		CreatedId int    `db:"created_id"`
		CreatedAt int    `db:"created_at"`
	}

	fileModel struct {
		table string
	}
)

// Insert 插入数据
func (m *fileModel) Insert(file *File) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?)", m.table, insertFileStr)
	result, err := db.Exec(sqlStr, file.Type, file.Url, file.Mime, file.Size, file.Filename, file.CreatedId)
	if err != nil {
		slog.Error("insert file err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	file.Id = int(lastInsertId)
	return nil
}

package model

import (
	"fmt"
	"golang.org/x/exp/slog"
)

var (
	UserModel = newUserModel()

	inertStr = "`username`,`password`,`phone_number`,`email`,`name`,`nick_name`,`avatar`,`create_id`,`update_id`"
)

type User struct {
	Model
	Username    string `db:"username"`
	Password    string `db:"password"`
	PhoneNumber string `db:"phone_number"`
	Email       string `db:"email"`
	Name        string `db:"name"`
	NickName    string `db:"nick_name"`
	Avatar      string `db:"avatar"`
}

type userModel struct {
	table string
}

func newUserModel() *userModel {
	return &userModel{
		table: "user",
	}
}

// Find 根据主键id单条查询
func (m *userModel) Find(id uint) (*User, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ?", m.table)
	user := new(User)
	if err := db.Get(user, sqlStr, id); err != nil {
		slog.Error("find user err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return user, nil
}

// Select 查询所有数据
func (m *userModel) Select() ([]*User, error) {
	sqlStr := fmt.Sprintf("select * from %s", m.table)
	var users []*User
	if err := db.Get(&users, sqlStr); err != nil {
		slog.Error("select user err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return users, nil
}

// SelectPage 分页查询
func (m *userModel) SelectPage(pageIndex, pageSize int) ([]*User, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	sqlStr := fmt.Sprintf("select * from %s limit %d,%d", m.table, (pageIndex-1)*pageSize, pageSize)
	var users []*User
	if err := db.Get(&users, sqlStr); err != nil {
		slog.Error("select page user err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return users, nil
}

// Insert 插入数据
func (m userModel) Insert(user *User) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?)", m.table, inertStr)
	result, err := db.Exec(sqlStr, user.Username, user.Password, user.PhoneNumber, user.Email, user.Name, user.NickName, user.Avatar, user.CreateId, user.UpdateId)
	if err != nil {
		slog.Error("insert user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	user.Id = uint(lastInsertId)
	return nil
}

// Update 更新数据
func (m userModel) Update(user *User) error {
	sqlStr := fmt.Sprintf("update %s set `phone_number` = ? , `email` = ? , `name`= ? , `nick_name`= ? , `acatar`= ? , `update_id` = ? where `id` = %d", m.table, user.Id)
	_, err := db.Exec(sqlStr, user.PhoneNumber, user.Email, user.Name, user.NickName, user.Avatar, user.UpdateId)
	if err != nil {
		slog.Error("update user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

// Delete 更新数据
func (m userModel) Delete(id uint) error {
	sqlStr := fmt.Sprintf("delete table %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("update user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

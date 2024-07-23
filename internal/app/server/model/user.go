package model

import (
	"fmt"
	"golang.org/x/exp/slog"
)

const (
	UserStatusEnable = iota + 1
	UserStatusDisable
)

var (
	UserModel = newUserModel()

	inertStr = "`username`,`password`,`phone_number`,`email`,`name`,`nick_name`,`avatar`,`state`,`create_id`,`update_id`"

	UserStatusDesc = map[int]string{
		UserStatusEnable:  "启用",
		UserStatusDisable: "禁用",
	}
)

type (
	User struct {
		Model
		Username    string `db:"username"`
		Password    string `db:"password"`
		PhoneNumber string `db:"phone_number"`
		Email       string `db:"email"`
		Name        string `db:"name"`
		NickName    string `db:"nick_name"`
		Avatar      string `db:"avatar"`
		State       int    `db:"state"`
		CreatedId   uint   `db:"created_id"`
		UpdatedId   uint   `db:"updated_id"`
	}

	userModel struct {
		table string
	}

	UserCond struct {
		Id          uint
		PhoneNumber string
		Name        string
		State       int
	}
)

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

// GetUsersByCond 根据条件分页获取用户
func (m *userModel) GetUsersByCond(cond *UserCond, pageIndex, pageSize int) ([]*User, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	sqlCond, args := m.buildUserCond(cond)
	sqlStr := fmt.Sprintf("select * from %s where 1 = 1 %s limit %d,%d", m.table, sqlCond, (pageIndex-1)*pageSize, pageSize)
	var users []*User
	if err := db.Select(&users, sqlStr, args...); err != nil {
		slog.Error("get users error ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return users, nil
}

// GetUserCountCond 根据条件获取用户数量
func (m *userModel) GetUserCountCond(cond *UserCond) (count int, err error) {
	sqlCond, args := m.buildUserCond(cond)
	sqlStr := fmt.Sprintf("select count(*) n from %s where 1 = 1  %s", m.table, sqlCond)
	if err = db.Get(&count, sqlStr, args...); err != nil {
		slog.Error("get user count error ", "sql", sqlStr, "err ", err.Error())
		return
	}
	return
}

// Insert 插入数据
func (m *userModel) Insert(user *User) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?,?)", m.table, inertStr)
	result, err := db.Exec(sqlStr, user.Username, user.Password, user.PhoneNumber, user.Email, user.Name, user.NickName, user.Avatar, user.State, user.CreatedId, user.UpdatedId)
	if err != nil {
		slog.Error("insert user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	user.Id = uint(lastInsertId)
	return nil
}

// Update 更新数据
func (m *userModel) Update(user *User) error {
	sqlStr := fmt.Sprintf("update %s set `phone_number` = ? , `email` = ? , `name`= ? , `nick_name`= ? , `acatar`= ? , `update_id` = ? where `id` = %d", m.table, user.Id)
	_, err := db.Exec(sqlStr, user.PhoneNumber, user.Email, user.Name, user.NickName, user.Avatar, user.UpdatedId)
	if err != nil {
		slog.Error("update user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

// Delete 更新数据
func (m *userModel) Delete(id uint) error {
	sqlStr := fmt.Sprintf("delete table %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("update user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *userModel) buildUserCond(cond *UserCond) (sqlCond string, args []interface{}) {
	if cond == nil {
		return
	}

	if cond.Id == 0 {
		sqlCond += "and id = ?"
		args = append(args, cond.Id)
	}

	if cond.PhoneNumber != "" {
		sqlCond += " and phone_number = ?"
		args = append(args, cond.PhoneNumber)
	}

	if cond.Name != "" {
		sqlCond += " and name = ?"
		args = append(args, cond.Name)
	}

	if cond.State > 0 {
		sqlCond += " and state = ?"
		args = append(args, cond.State)
	}
	return
}

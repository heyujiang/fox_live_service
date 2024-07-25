package model

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"time"
)

const (
	UserStatusEnable = iota + 1
	UserStatusDisable

	inertStr = "`username`,`password`,`phone_number`,`email`,`name`,`nick_name`,`avatar`,`state`,`created_id`,`updated_id`"
)

var (
	UserModel = newUserModel()

	UserStatusDesc = map[int]string{
		UserStatusEnable:  "启用",
		UserStatusDisable: "禁用",
	}
)

type (
	User struct {
		Id          int       `db:"id"`
		Username    string    `db:"username"`
		Password    string    `db:"password"`
		PhoneNumber string    `db:"phone_number"`
		Email       string    `db:"email"`
		Name        string    `db:"name"`
		NickName    string    `db:"nick_name"`
		Avatar      string    `db:"avatar"`
		State       int       `db:"state"`
		CreatedId   int       `db:"created_id"`
		UpdatedId   int       `db:"updated_id"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	userModel struct {
		table string
	}

	UserCond struct {
		Id          int
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

// Insert 插入数据
func (m *userModel) Insert(user *User) error {
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?,?)", m.table, inertStr)
	result, err := db.Exec(sqlStr, user.Username, user.Password, user.PhoneNumber, user.Email, user.Name, user.NickName, user.Avatar, user.State, user.CreatedId, user.UpdatedId)
	if err != nil {
		slog.Error("insert user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	user.Id = int(lastInsertId)
	return nil
}

// Delete 更新数据
func (m *userModel) Delete(id int) error {
	sqlStr := fmt.Sprintf("delete table %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("update user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

// Update 更新数据
func (m *userModel) Update(user *User) error {
	sqlStr := fmt.Sprintf("update %s set `email` = ? , `nick_name`= ? , `avatar`= ? , `updated_id` = ? where `id` = %d", m.table, user.Id)
	_, err := db.Exec(sqlStr, user.Email, user.NickName, user.Avatar, user.UpdatedId)
	if err != nil {
		slog.Error("update user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

// Find 根据主键id单条查询
func (m *userModel) Find(id int) (*User, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ?", m.table)
	user := new(User)
	if err := db.Get(user, sqlStr, id); err != nil {
		slog.Error("find user err ", "sql", sqlStr, "id", id, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
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

// FindByUsername 根据username查询用户信息
func (m *userModel) FindByUsername(username string) (*User, error) {
	sqlStr := fmt.Sprintf("select * from %s where `username` = ?", m.table)
	user := new(User)
	if err := db.Get(user, sqlStr, username); err != nil {
		slog.Error("find user by username error", "sql", sqlStr, "username", username, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return user, nil
}

// FindByPhoneNumber 根据 phoneNumber 查询用户信息
func (m *userModel) FindByPhoneNumber(phoneNumber string) (*User, error) {
	sqlStr := fmt.Sprintf("select * from %s where `phone_number` = ?", m.table)
	user := new(User)
	if err := db.Get(user, sqlStr, phoneNumber); err != nil {
		slog.Error("find user by phone number error", "sql", sqlStr, "phone_number", phoneNumber, "err ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotRecord
		}
		return nil, err
	}
	return user, nil
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

func (m *userModel) buildUserCond(cond *UserCond) (sqlCond string, args []interface{}) {
	if cond == nil {
		return
	}

	if cond.Id > 0 {
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

	if _, ok := UserStatusDesc[cond.State]; ok {
		sqlCond += " and state = ?"
		args = append(args, cond.State)
	}

	return
}

func (m *userModel) UpdateState(id, state, uid int) error {
	sqlStr := fmt.Sprintf("update %s set `state` = ? , `updated_id` = ?  where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, state, uid)
	if err != nil {
		slog.Error("update user state err ", "sql", sqlStr, "id", id, "state", state, "err ", err.Error())
		return err
	}
	return nil
}

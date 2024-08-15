package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"
	"time"
)

const (
	UserStatusEnable = iota + 1
	UserStatusDisable

	inertStr = "`username`,`password`,`phone_number`,`email`,`name`,`nick_name`,`avatar`,`state`,`role_ids`,`dept_id`,`created_id`,`updated_id`"
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
		RoleIds     string    `db:"role_ids"`
		DeptId      int       `db:"dept_id"`
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
		Username    string
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
	sqlStr := fmt.Sprintf("insert into %s (%s) values (?,?,?,?,?,?,?,?,?,?,?,?)", m.table, inertStr)
	result, err := db.Exec(sqlStr, user.Username, user.Password, user.PhoneNumber, user.Email, user.Name, user.NickName, user.Avatar, user.State, user.RoleIds, user.DeptId, user.CreatedId, user.UpdatedId)
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
	sqlStr := fmt.Sprintf("delete from %s where `id` = ? ", m.table)
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		slog.Error("update user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

// Update 更新数据
func (m *userModel) Update(user *User) error {
	sqlStr := fmt.Sprintf("update %s set `email` = ? ,`name` = ?, `nick_name`= ? , `avatar`= ? , `role_ids` = ?, `dept_id` = ?, `updated_id` = ? where `id` = %d", m.table, user.Id)
	_, err := db.Exec(sqlStr, user.Email, user.Name, user.NickName, user.Avatar, user.RoleIds, user.DeptId, user.UpdatedId)
	if err != nil {
		slog.Error("update user err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

// Find 根据主键id单条查询
func (m *userModel) Find(id int) (*User, error) {
	sqlStr := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
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
	sqlStr := fmt.Sprintf("select * from %s where `state` = ?", m.table)
	var users []*User
	if err := db.Select(&users, sqlStr, UserStatusEnable); err != nil {
		slog.Error("select user err ", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return users, nil
}

// FindByUsername 根据username查询用户信息
func (m *userModel) FindByUsername(username string) (*User, error) {
	sqlStr := fmt.Sprintf("select * from %s where `username` = ? limit 1", m.table)
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
	sqlStr := fmt.Sprintf("select * from %s where `phone_number` = ? limit 1", m.table)
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
	sqlStr := fmt.Sprintf("select count(*) n from %s where 1 = 1 %s", m.table, sqlCond)
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

	if cond.Username != "" {
		sqlCond += " and username = ?"
		args = append(args, cond.Username)
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

func (m *userModel) SelectByEnable() ([]*User, error) {
	var users []*User
	sqlStr := fmt.Sprintf("select * from %s where `state` = ?", m.table)
	if err := db.Select(&users, sqlStr, UserStatusEnable); err != nil {
		slog.Error("find user options error", "sql", sqlStr, "err ", err.Error())
		return nil, err
	}
	return users, nil
}

func (m *userModel) SelectByIds(ids []int) ([]*User, error) {
	var users []*User
	sqlStr := fmt.Sprintf("select * from %s where `state` = ? and id in (?) ", m.table)
	query1, args, err := sqlx.In(sqlStr, UserStatusEnable, ids)
	if err != nil {
		slog.Error("batch select user bu uids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}
	slog.Info(query1, "args", args)
	if err := db.Select(&users, query1, args...); err != nil {
		slog.Error("batch select user bu uids error", "sql", sqlStr, "ids", ids, "err ", err.Error())
		return nil, err
	}
	return users, nil
}

func (m *userModel) UpdateAvatar(id int, avatar string) error {
	sqlStr := fmt.Sprintf("update %s set `avatar` = ? , `updated_id` = ?  where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, avatar, id)
	if err != nil {
		slog.Error("update user avatar err ", "sql", sqlStr, "id", id, "avatar", avatar, "err ", err.Error())
		return err
	}
	return nil
}

func (m *userModel) UpdateBasic(user *User) error {
	sqlStr := fmt.Sprintf("update %s set `email` = ?  , `updated_id` = ?  where `id` = %d", m.table, user.Id)
	_, err := db.Exec(sqlStr, user.Email, user.UpdatedId)
	if err != nil {
		slog.Error("update user basic err ", "sql", sqlStr, "err ", err.Error())
		return err
	}
	return nil
}

func (m *userModel) UpdatePassword(id int, password string) error {
	sqlStr := fmt.Sprintf("update %s set `password` = ? , `updated_id` = ?  where `id` = %d", m.table, id)
	_, err := db.Exec(sqlStr, password, id)
	if err != nil {
		slog.Error("update user password err ", "sql", sqlStr, "id", id, "err ", err.Error())
		return err
	}
	return nil
}

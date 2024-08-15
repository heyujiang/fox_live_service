package model

import (
	"errors"
	"fmt"
	"fox_live_service/config/global"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type (
	UserCountItem struct {
		UserId int `db:"user_id"`
		Count  int `db:"count"`
	}
)

var db *sqlx.DB
var dbType = "mysql"

var ErrNotRecord = errors.New("record is not exist")

func init() {
	db = sqlx.MustConnect(dbType, global.Config.GetString("Db.Mysql.DSN"))
	db.SetMaxOpenConns(global.Config.GetInt("Db.Mysql.MaxOpenConn"))
	db.SetMaxIdleConns(global.Config.GetInt("Db.Mysql.MaxIdleConn"))
}

type Model struct {
	Id        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func transaction(tx *sqlx.Tx, fn func() error) (err error) {
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
	err = fn()
	return err
}

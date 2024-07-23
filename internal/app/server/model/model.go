package model

import (
	"fox_live_service/config/global"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var db *sqlx.DB
var dbType = "mysql"

func init() {
	db = sqlx.MustConnect(dbType, global.Config.GetString("Db.Mysql.DSN"))
	db.SetMaxOpenConns(global.Config.GetInt("Db.Mysql.MaxOpenConn"))
	db.SetMaxIdleConns(global.Config.GetInt("Db.Mysql.MaxIdleConn"))
}

type Model struct {
	Id        uint      `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

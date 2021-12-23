package mariadb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"jwtToken/cfg"
	"os"
)

var db *sql.DB

func InitDB() {
	//db, _ = sql.Open("mysql", "root:123@tcp(localhost:3306)/ljhdb")
	cfg := cfg.Cfg.Db
	db, _ = sql.Open(cfg.Driver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName))

	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
}

// DBConn : 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

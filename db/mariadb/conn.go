package mariadb

import (
	"database/sql"
	"fmt"
	//config
	"jwtToken/config"

	//driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func NewDatabase(cnf *config.Config) error {
	var err error

	db, err = sql.Open(cnf.Db.Type, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cnf.Db.User, cnf.Db.Password, cnf.Db.Host, cnf.Db.Port, cnf.Db.DatabaseName))

	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	// Max idle connections
	db.SetMaxOpenConns(cnf.Db.MaxIdleConns)
	// Max open connections
	db.SetMaxIdleConns(cnf.Db.MaxOpenConns)

	return nil
}

// DBConn : 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

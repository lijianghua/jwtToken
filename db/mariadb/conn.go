package mariadb

import (
	"database/sql"
	"fmt"
	//_ "github.com/golang-migrate/migrate/v4/source/file"
	//config
	"jwtToken/config"

	//driver
	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase(dbCnf *config.DatabaseConfig) (*sql.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbCnf.User, dbCnf.Password, dbCnf.Host, dbCnf.Port, dbCnf.DatabaseName)

	db, err := sql.Open(dbCnf.Type, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Max idle connections
	db.SetMaxOpenConns(dbCnf.MaxIdleConns)
	// Max open connections
	db.SetMaxIdleConns(dbCnf.MaxOpenConns)

	//mCfg := &migratemysql.Config{
	//	DatabaseName: dbCnf.DatabaseName,
	//}
	//driver, err := migratemysql.WithInstance(db, mCfg)
	//if err != nil {
	//	return nil, err
	//}
	//m, err := migrate.NewWithDatabaseInstance("file://"+dbCnf.MigrationDir, "mysql", driver)
	//if err != nil {
	//	return nil, err
	//}
	//err = m.Up()
	//if err != nil && err != migrate.ErrNoChange {
	//	return nil, err
	//}
	return db, nil
}

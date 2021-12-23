package db

import (
	"jwtToken/db/mariadb"
)

type TableUser struct {
	UserID   int64
	Username string
	Password string
}

func NewUser(username, password string) error {
	_, err := mariadb.DBConn().Exec("insert into tbl_user(user_name,user_pwd) values(?,?)", username, password)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(username string) (*TableUser, error) {
	user := TableUser{}

	err := mariadb.DBConn().QueryRow("select id,user_name,user_pwd from tbl_user where user_name = ? limit 1", username).Scan(&user.UserID, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

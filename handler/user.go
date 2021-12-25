package handler

import (
	"jwtToken/db"
	"jwtToken/util"
)

type User struct {
	UserID   int64
	Username string
	Password string
}

func (m *User) Load(username string) error {
	user, err := db.GetUser(username)
	if err != nil {
		return err
	}
	m.UserID = user.UserID
	m.Username = user.Username

	return nil
}
func (m *User) IsValid() bool {
	//check parameter
	if len(m.Username) < 3 || len(m.Password) < 3 {
		return false
	}
	return true
}

func (m *User) Verify() bool {

	if valid := m.IsValid(); !valid {
		return false
	}
	//query user info
	user, err := db.GetUser(m.Username)
	if err != nil {
		return false
	}

	//verify password
	suc := util.VerifyPass(user.Password, m.Password)
	if !suc {
		return false
	}

	return true

}

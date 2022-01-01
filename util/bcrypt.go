package util

import (
	"golang.org/x/crypto/bcrypt"
)

//Bcrypt hash password
func HashPass(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//Bcrypt verify password
func VerifyPass(encodedPwd string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPwd), []byte(pwd))
	return err == nil
}

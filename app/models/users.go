package models

import (
	"encoding/base64"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserName    string `gorm:"column:user_name;size:50"`
	UserHashPwd string `gorm:"column:user_password;size:100"`
	UserEMail   string `gorm:"column:user_email;size:100"`
}

func (user *User) GetPasswordHash() (hash []byte, err error) {
	hash, err = base64.StdEncoding.DecodeString(user.UserHashPwd)
	return
}

func (user *User) SetPasswordHash(hash []byte) {
	user.UserHashPwd = base64.StdEncoding.EncodeToString(hash)
}

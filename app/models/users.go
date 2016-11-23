package models

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
)

const (
	GroupPackager = 0   // this user view only himself builds
	GroupPusher   = 1   // this user view and may be a push all packages
	GroupAdmin    = 100 // thos user may edit users and push all packages
)

type User struct {
	gorm.Model
	UserName    string `gorm:"column:user_name;size:50"`
	UserHashPwd string `gorm:"column:user_password;size:100"`
	UserEMail   string `gorm:"column:user_email;size:100"`
	UserGroup   int    `gorm:"column:group"`
}

type Group int

func (user *User) GetPasswordHash() (hash []byte, err error) {
	hash, err = base64.StdEncoding.DecodeString(user.UserHashPwd)
	return
}

func (user *User) SetPasswordHash(hash []byte) {
	user.UserHashPwd = base64.StdEncoding.EncodeToString(hash)
}

func (user *User) GeneratePasswordHash(pwd string) (hash []byte) {
	hash, _ = bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return
}

package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName    string `gorm:"column:user_name;size:50"`
	UserHashPwd string `gorm:"column:user_password;size:100"`
	UserEMail   string `gorm:"column:user_email;size:100"`
}

package models

import "github.com/jinzhu/gorm"

type Owner struct {
	gorm.Model
	OwnerName string `gorm:"column:owner_name;size:50"`
}

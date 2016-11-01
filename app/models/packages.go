package models

import "github.com/jinzhu/gorm"

type Package struct {
	gorm.Model
	PkgName  string `gorm:"column:pkg_name;unique_index"`
	PkgRepo  Repo   `gorm:"column:pkg_repo"`
	PkgOwner Owner  `gorm:"column:pkg_owner"`
}

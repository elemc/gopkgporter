package models

import "github.com/jinzhu/gorm"

type Package struct {
	gorm.Model
	PkgName    string `gorm:"column:pkg_name;unique_index"`
	PkgRepo    Repo
	PkgRepoID  uint `gorm:"column:pkg_repo_id"`
	PkgOwner   Owner
	PkgOwnerID uint `gorm:"column:pkg_owner_id"`
}

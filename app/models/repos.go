package models

import "github.com/jinzhu/gorm"

type RepoType struct {
	gorm.Model
	RTName string `gorm:"column:rt_name"`
}

type Repo struct {
	gorm.Model
	RepoName string `gorm:"column:repo_name"`
}

// String returns repository string view / name
func (repo *Repo) String() string {
	return repo.RepoName
}

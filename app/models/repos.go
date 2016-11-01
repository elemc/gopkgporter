package models

import "github.com/jinzhu/gorm"

type RepoType struct {
	gorm.Model
	RTName string `gorm:"column:rt_name;unique_index"`
}

type Repo struct {
	gorm.Model
	RepoName string `gorm:"column:repo_name;unique_index"`
}

// String returns repository string view / name
func (repo *Repo) String() string {
	return repo.RepoName
}

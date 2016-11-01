package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BuildedPackage struct {
	gorm.Model
	BuildID        uint    `gorm:"column:build_id"`
	BuildPackage   Package `gorm:"column:build_pkg"`
	BuildPackageID uint
	Version        string    `gorm:"column:version;size:50"`
	Release        string    `gorm:"column:release;size:25"`
	Epoch          string    `gorm:"column:epoch"`
	CompletionTime time.Time `gorm:"column:completion_time"`
	TaskID         uint      `gorm:"column:task_id"`
	Owner          Owner     `gorm:"column:owner"`
	OwnerID        uint
	Pushed         bool `gorm:"column:pushed"`
	WaitToTime     bool `gorm:"column:wait_to_time"`
	PushUser       User `gorm:"column:push_user"`
	PushUserID     uint
	PushTime       time.Time `gorm:"column:push_time"`
	PushRepoType   RepoType  `gorm:"column:push_repo_type"`
	PushRepoTypeID uint
	BlockedToPush  bool   `gorm:"column:is_blocked_to_push"`
	TagName        string `gorm:"column:tag_name;size:25"`
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Log struct {
	gorm.Model
	Timestamp time.Time `gorm:"column:timestamp"`
	Package   Package
	PackageID int
	Action    string `gorm:"column:action;size:100"`
	User      User
	UserID    int
	Type      string `gorm:"column:type;size:10"`
}

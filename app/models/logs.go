package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Log model for logs
type Log struct {
	gorm.Model
	Timestamp time.Time `gorm:"column:timestamp"`
	Package   Package
	PackageID uint
	Action    string `gorm:"column:action;size:100"`
	User      User
	UserID    uint   `gorm:"column:user_id"`
	Type      string `gorm:"column:type;size:10"`
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Log struct {
	gorm.Model
	Timestamp time.Time `gorm:"column:timestamp"`
	Package   Package   `gorm:"column:package"`
	Action    string    `gorm:"column:action;size:25"`
	User      User      `gorm:"column:user"`
}

package controllers

import (
	"gopkgporter/app/models"

	"github.com/revel/revel"
)

const (
	timeFormat = "2006-01-02 15:04:05 -0700 MST"
)

// App is a root struct for web application
type App struct {
	*revel.Controller
}

// Index function returns log list
func (c App) Index() revel.Result {
	//var users []models.User
	//dbgorm.Db.Find(&users)

	var logs []models.Log
	logsQuery := dbgorm.Db.Order("Timestamp DESC", true).Find(&logs).Limit(25)
	err := logsQuery.Error
	if err != nil {
		return c.RenderError(err)
	}

	var newLogs []models.Log
	for _, log := range logs {
		dbgorm.Db.Model(&log).Related(&log.User)
		dbgorm.Db.Model(&log).Related(&log.Package)
		newLogs = append(newLogs, log)
	}
	// show logs
	return c.Render(newLogs, timeFormat)
}

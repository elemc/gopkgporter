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

	currentUser := connected(c.RenderArgs, c.Session)

	// show logs
	return c.Render(newLogs, timeFormat, currentUser)
}

func connected(args map[string]interface{}, session revel.Session) *models.User {
	if args["user"] != nil {
		return args["user"].(*models.User)
	}
	if username, ok := session["user"]; ok {
		return getUser(username)
	}
	return nil
}

func getUser(username string) *models.User {
	var users []models.User
	ctx := dbgorm.Db.Find(&users, "user_name=?", username)
	if ctx.Error != nil {
		revel.ERROR.Printf("Error in getUser: %s", ctx.Error)
		return nil
	}
	if len(users) == 0 {
		revel.WARN.Printf("User with username %s not found in database", username)
		return nil
	}

	return &users[0]
}

package controllers

import (
	"gopkgporter/app/models"

	"github.com/revel/revel"
)

const (
	dontPerm = "You don't have permissions for update this package!"
)

func init() {
	revel.OnAppStart(InitDB)
	// revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	// revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	// revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
	// revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	// revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
	revel.TemplateFuncs["user_is_auth"] = func(s revel.Session) bool {
		_, ok := s["user"]
		return ok
	}
	revel.TemplateFuncs["auth_user_name"] = func(s revel.Session) string {
		if username, ok := s["user"]; ok {
			return username
		}
		return ""
	}
	revel.TemplateFuncs["user_is_nil"] = func(user *models.User) bool {
		return user == nil
	}
	revel.TemplateFuncs["user_is_admin"] = func(user *models.User) bool {
		if user == nil {
			return false
		}
		return user.UserGroup == models.GroupAdmin
	}
	revel.TemplateFuncs["user_is_pusher"] = func(user *models.User) bool {
		if user == nil {
			return false
		}
		return (user.UserGroup == models.GroupAdmin) || (user.UserGroup == models.GroupPusher)
	}
}

package controllers

import (
	"github.com/elemc/gopkgporter/app/models"
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
	revel.TemplateFuncs["auth_user_name"] = func(user *models.User) string {
		if user != nil {
			return user.UserName
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

	revel.TemplateFuncs["get_group_name"] = func(user models.User) string {
		var groupName string
		switch user.UserGroup {
		case models.GroupAdmin:
			groupName = "Admin"
		case models.GroupPusher:
			groupName = "Pusher"
		case models.GroupPackager:
			groupName = "Packager"
		default:
			groupName = "Unknown"
		}
		return groupName
	}
	revel.TemplateFuncs["is_current_user"] = func(user models.User, currentUser *models.User) bool {
		return user.ID == currentUser.ID
	}
}

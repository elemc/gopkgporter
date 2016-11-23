package controllers

import "github.com/revel/revel"

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
}

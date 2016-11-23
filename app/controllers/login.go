package controllers

import (
	"testesia/app/routes"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	*revel.Controller
}

func (c Login) Index() revel.Result {
	user := connected(c.RenderArgs, c.Session)
	if user != nil {
		c.Flash.Error("User %s is logged on", user.UserName)
		return c.Redirect("/")
	}
	return c.Render()
}

func (c Login) Login(username, password string, remember bool) revel.Result {
	user := getUser(username)
	if user != nil {
		hash, err := user.GetPasswordHash()
		if err != nil {
			revel.ERROR.Printf("Error per login: %s", err)
			c.Flash.Error("Error in login", err)
			return c.Redirect("/")
		}
		err = bcrypt.CompareHashAndPassword(hash, []byte(password))
		if err == nil {
			c.Session["user"] = username
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, %s", user.UserName)
			c.RenderArgs["user"] = user
			return c.Redirect("/")
		}
	}
	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.App.Index())
}

func (c Login) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

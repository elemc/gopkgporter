package controllers

import (
	"github.com/elemc/gopkgporter/app/models"
	"github.com/elemc/gopkgporter/app/routes"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

// Login controller /login
type Login struct {
	*revel.Controller
}

// Index function for index /login
func (c Login) Index() revel.Result {
	user := connected(c.RenderArgs, c.Session)
	if user != nil {
		c.Flash.Error("User %s is logged on", user.UserName)
		return c.Redirect("/")
	}
	return c.Render()
}

// Login function for login
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
			session := models.CreateSession()
			session.SessionUser = user
			session.SessionUserID = user.ID
			ctx := dbgorm.Db.Create(session)
			if ctx.Error != nil {
				revel.ERROR.Printf("Error in create session: %s", ctx.Error)
				return c.Redirect(routes.App.Index())
			}
			ctx = dbgorm.Db.Save(session)
			if ctx.Error != nil {
				revel.ERROR.Printf("Error in save session: %s", ctx.Error)
				return c.Redirect(routes.App.Index())
			}

			c.Session["user"] = session.SessionID
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

// Logout function for /logout
func (c Login) Logout() revel.Result {
	if sessionID, ok := c.Session["user"]; ok {
		session := getSession(sessionID)
		ctx := dbgorm.Db.Delete(session, "session_id=?", sessionID)
		if ctx.Error != nil {
			revel.ERROR.Printf("Error in delete session: %s", ctx.Error)
			return c.Redirect(routes.App.Index())
		}
	}
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

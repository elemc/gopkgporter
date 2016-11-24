package controllers

import (
	"gopkgporter/app/models"

	"github.com/revel/revel"
)

type Users struct {
	*revel.Controller
}

func (c *Users) Index() revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)
	if currentUser == nil || currentUser.UserGroup < models.GroupAdmin {
		return c.Render(currentUser)
	}

	var users []models.User
	ctx := dbgorm.Db.Find(&users)
	if ctx.Error != nil {
		return c.RenderError(ctx.Error)
	}

	return c.Render(users, currentUser)
}

func (c *Users) Edit(id uint) revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)
	// if currentUser == nil || (currentUser.UserGroup < models.GroupAdmin && currentUser.ID != id) {
	// 	return c.Render(currentUser)
	// }

	var user models.User
	dbgorm.Db.First(&user, id)
	return c.Render(user, currentUser)
}

func (c *Users) Save(id uint, username, password, email string, group int) revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)
	if currentUser == nil {
		c.Flash.Error(dontPerm)
		return c.Redirect("/users")
	}

	var user models.User
	isNew := false
	ctx := dbgorm.Db.First(&user, id)
	if ctx.Error != nil {
		isNew = true
	}
	if password != "" {
		user.SetPasswordHash(user.GeneratePasswordHash(password))
	}
	user.UserEMail = email

	if currentUser.UserGroup >= models.GroupAdmin {
		user.UserName = username
		user.UserGroup = group
	}

	if isNew {
		ctx = dbgorm.Db.Create(&user)
	} else {
		ctx = dbgorm.Db.Save(&user)
	}
	if ctx.Error != nil {
		return c.RenderError(ctx.Error)
	}

	return c.Redirect("/users")
}

func (c *Users) Delete(id uint) revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)
	if currentUser == nil || currentUser.UserGroup < models.GroupAdmin {
		c.Flash.Error(dontPerm)
		return c.Redirect("/users")
	}

	var user models.User
	ctx := dbgorm.Db.First(&user, id)
	if ctx.Error != nil {
		return c.RenderError(ctx.Error)
	}
	ctx = dbgorm.Db.Delete(&user, id)
	if ctx.Error != nil {
		return c.RenderError(ctx.Error)
	}
	return c.Redirect("/users")
}

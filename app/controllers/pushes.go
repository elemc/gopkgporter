package controllers

import (
	"fmt"

	"github.com/elemc/gopkgporter/app/models"
	"github.com/revel/revel"
)

// Pushes controller for /pushes
type Pushes struct {
	*revel.Controller
}

// Index function for index /pushes page
func (c Pushes) Index() revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)

	var pushes []models.PackagesToPush
	q := dbgorm.Db.Find(&pushes, "done='false'")
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	return c.Render(pushes, currentUser)
}

// Edit function for save push object
func (c Pushes) Edit(id int) revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)
	if currentUser == nil || currentUser.UserGroup < models.GroupAdmin {
		return c.RenderError(fmt.Errorf(dontPerm))
	}

	var push models.PackagesToPush
	q := dbgorm.Db.Find(&push, id)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	push.Version = c.Params.Get("Version")
	push.Repository = c.Params.Get("Repository")
	push.Branch = c.Params.Get("Branch")
	push.Distributive = c.Params.Get("Distributive")

	q = dbgorm.Db.Save(&push)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	c.Flash.Success("Build for push saved successfully")
	return c.Redirect("/pushes")

}

// Delete function for POST method deleted push for selected id
func (c Pushes) Delete(id int) revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)
	if currentUser == nil || currentUser.UserGroup < models.GroupAdmin {
		return c.RenderError(fmt.Errorf(dontPerm))
	}
	var push models.PackagesToPush
	q := dbgorm.Db.Find(&push, id)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}
	q = dbgorm.Db.Delete(&push)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	c.Flash.Success("Build for push removed successfully.")
	return c.Redirect("/pushes")
}

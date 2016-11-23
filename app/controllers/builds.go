package controllers

import (
	"fmt"
	"gopkgporter/app/models"
	"time"

	"github.com/revel/revel"
)

// Builds controller for work with builds
type Builds struct {
	*revel.Controller
}

// Index function returns builded but not pushed packages
func (c Builds) Index() revel.Result {
	var rbuilds []models.BuildedPackage
	var builds []models.BuildedPackage

	where := "pushed='false' AND is_blocked_to_push='false'"

	currentUser := connected(c.RenderArgs, c.Session)
	if currentUser != nil {
		if currentUser.UserGroup < models.GroupPusher {
			var owner models.Owner
			ctx := dbgorm.Db.First(&owner, "owner_name=?", currentUser.UserName)
			if ctx.Error == nil {
				where = fmt.Sprintf("%s AND owner_id='%d'", where, owner.ID)
			} else {
				where = fmt.Sprintf("%s AND owner_id=-1", where)
				c.Flash.Error("Not found owner with name: %s", currentUser.UserName)
			}
		}
	}
	dbgorm.Db.Order("build_id DESC", true).Find(&rbuilds, where)

	for _, build := range rbuilds {
		dbgorm.Db.Model(&build).Related(&build.BuildPackage, "BuildPackage")
		dbgorm.Db.Model(&build).Related(&build.Owner, "Owner")
		dbgorm.Db.Model(&build).Related(&build.User, "PushUser")
		dbgorm.Db.Model(&build).Related(&build.PushRepoType, "PushRepoType")
		builds = append(builds, build)
	}

	return c.Render(builds, currentUser)
}

// Get function returns build information
func (c Builds) Get(id int) revel.Result {
	return c.Redirect("/")
}

// CancelBuild function canceled build to push
func (c Builds) CancelBuild(id int) revel.Result {
	var builds []models.BuildedPackage
	q := dbgorm.Db.Find(&builds, fmt.Sprintf("build_id=%d", id))
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	currentUser := connected(c.RenderArgs, c.Session)

	for _, build := range builds {
		dbgorm.Db.Model(&build).Related(&build.BuildPackage, "BuildPackage")
		dbgorm.Db.Model(&build).Related(&build.Owner, "Owner")
		dbgorm.Db.Model(&build).Related(&build.User, "PushUser")
		dbgorm.Db.Model(&build).Related(&build.PushRepoType, "PushRepoType")

		build.BlockedToPush = true
		build.PushTime = time.Now()

		dbgorm.Db.Find(&build.User, currentUser.ID)

		q = dbgorm.Db.Save(&build)
		if q.Error != nil {
			return c.RenderError(q.Error)
		}
	}

	return c.Redirect("/builds")
}

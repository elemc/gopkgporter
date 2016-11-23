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
	dbgorm.Db.Order("build_id DESC", true).Find(&rbuilds, "pushed='false' AND is_blocked_to_push='false'")

	for _, build := range rbuilds {
		dbgorm.Db.Model(&build).Related(&build.BuildPackage, "BuildPackage")
		dbgorm.Db.Model(&build).Related(&build.Owner, "Owner")
		dbgorm.Db.Model(&build).Related(&build.User, "PushUser")
		dbgorm.Db.Model(&build).Related(&build.PushRepoType, "PushRepoType")
		builds = append(builds, build)
	}

	return c.Render(builds)
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

	for _, build := range builds {
		build.BlockedToPush = true
		build.PushTime = time.Now()
		dbgorm.Db.Find(&build.User, 1)

		q = dbgorm.Db.Save(&build)
		if q.Error != nil {
			return c.RenderError(q.Error)
		}
	}

	return c.Redirect("/builds")
}

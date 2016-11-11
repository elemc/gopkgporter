package controllers

import (
	"gopkgporter/app/models"

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
	dbgorm.Db.Order("build_id DESC", true).Find(&rbuilds, "pushed='false'")

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

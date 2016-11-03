package controllers

import (
	"gopkgporter/app/models"

	"github.com/revel/revel"
)

type Builds struct {
	*revel.Controller
}

func (c Builds) Index() revel.Result {
	var rbuilds []models.BuildedPackage
	var builds []models.BuildedPackage
	dbgorm.Db.Find(&rbuilds, "pushed='false'")

	for _, build := range rbuilds {
		dbgorm.Db.Model(&build).Related(&build.BuildPackage, "BuildPackage")
		dbgorm.Db.Model(&build).Related(&build.Owner)
		dbgorm.Db.Model(&build).Related(&build.User, "PushUser")
		dbgorm.Db.Model(&build).Related(&build.PushRepoType, "PushRepoType")
		revel.INFO.Printf("%+v", build.BuildPackage)
		builds = append(builds, build)
	}

	return c.Render(builds)
}

func (c Builds) Get(id int) revel.Result {
	return c.Redirect("/")
}

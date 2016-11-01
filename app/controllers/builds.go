package controllers

import (
	"gopkgporter/app/models"

	"github.com/revel/revel"
)

type Builds struct {
	*revel.Controller
}

func (c Builds) Index() revel.Result {
	var builds []models.BuildedPackage
	dbgorm.Db.Find(&builds)

	return c.Render()
}

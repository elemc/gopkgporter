package controllers

import (
	"fmt"
	"gopkgporter/app/models"
	"strconv"

	"github.com/revel/revel"
)

type Packages struct {
	*revel.Controller
}

func (c Packages) Index() revel.Result {
	var pkgs []models.Package
	q := dbgorm.Db.Order("pkg_name ASC", true).Find(&pkgs)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}
	var newPkgs []models.Package
	for _, pkg := range pkgs {
		dbgorm.Db.Model(&pkg).Related(&pkg.PkgOwner, "PkgOwner")
		dbgorm.Db.Model(&pkg).Related(&pkg.PkgRepo, "PkgRepo")
		newPkgs = append(newPkgs, pkg)
	}
	pkgs = nil

	var ownerList []models.Owner

	q = dbgorm.Db.Order("owner_name ASC", true).Find(&ownerList)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	var repos []models.Repo
	q = dbgorm.Db.Order("repo_name ASC", true).Find(&repos)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	return c.Render(newPkgs, ownerList, repos)
}

func (c Packages) Edit(id int) revel.Result {
	var pkg models.Package
	q := dbgorm.Db.Find(&pkg, id)
	if q.Error != nil {
		c.RenderError(q.Error)
	}

	ownerID := c.Params.Get("OwnerID")
	if ownerID != "" {
		oid, err := strconv.Atoi(ownerID)
		var owner models.Owner
		if err == nil {
			q := dbgorm.Db.Find(&owner, oid)
			if q.Error != nil {
				return c.RenderError(q.Error)
			}
			pkg.PkgOwner = owner
		}
	}

	repoID := c.Params.Get("RepoID")
	if repoID != "" {
		oid, err := strconv.Atoi(repoID)
		var repo models.Repo
		if err == nil {
			q := dbgorm.Db.Find(&repo, oid)
			if q.Error != nil {
				return c.RenderError(q.Error)
			}
			pkg.PkgRepo = repo
		}
	}

	dbgorm.Db.Save(&pkg)

	return c.Redirect("/packages")
}

func (c Packages) Package(id int) revel.Result {
	var pkg models.Package
	q := dbgorm.Db.Find(&pkg, id)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	dbgorm.Db.Model(&pkg).Related(&pkg.PkgOwner, "PkgOwner")
	dbgorm.Db.Model(&pkg).Related(&pkg.PkgRepo, "PkgRepo")

	var ownerList []models.Owner

	q = dbgorm.Db.Order("owner_name ASC", true).Find(&ownerList)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	var repos []models.Repo
	q = dbgorm.Db.Order("repo_name ASC", true).Find(&repos)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	titleName := fmt.Sprintf("Package \"%s\"", pkg.PkgName)

	return c.Render(pkg, ownerList, repos, titleName)
}

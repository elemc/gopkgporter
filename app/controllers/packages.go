package controllers

import (
	"fmt"
	"strconv"

	"github.com/elemc/gopkgporter/app/models"
	"github.com/revel/revel"
)

// Packages controller
type Packages struct {
	*revel.Controller
}

// Index function returns list packages
func (c Packages) Index() revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)

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

	return c.Render(newPkgs, ownerList, repos, currentUser)
}

// Edit function edit owner or repository for specified package
func (c Packages) Edit(id int) revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)
	var pkg models.Package
	q := dbgorm.Db.Find(&pkg, id)
	if q.Error != nil {
		c.RenderError(q.Error)
	}

	if currentUser == nil {
		return c.RenderError(fmt.Errorf(dontPerm))
	} else if currentUser.UserGroup < models.GroupPusher {
		var owner models.Owner
		ctx := dbgorm.Db.First(&owner, "owner_name=?", currentUser.UserName)
		if ctx.Error != nil || owner.ID != pkg.PkgOwnerID {
			return c.RenderError(fmt.Errorf(dontPerm))
		}
	}

	ownerID := c.Params.Get("OwnerID")
	if ownerID != "" {
		oid, err := strconv.Atoi(ownerID)
		var owner models.Owner
		if err == nil {
			q = dbgorm.Db.Find(&owner, oid)
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
			q = dbgorm.Db.Find(&repo, oid)
			if q.Error != nil {
				return c.RenderError(q.Error)
			}
			pkg.PkgRepo = repo
		}
	}

	q = dbgorm.Db.Save(&pkg)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}
	c.Flash.Success("Package \"%s\" successfully saved.", pkg.PkgName)

	returnPage := "/packages"
	rPage := c.Params.Get("return_page")
	if rPage != "" {
		returnPage = rPage
	}

	return c.Redirect(returnPage)
}

// Package function returns page with one specified package with id
func (c Packages) Package(id int) revel.Result {
	currentUser := connected(c.RenderArgs, c.Session)
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

	returnPage := "/packages"
	rPage := c.Params.Get("return_page")
	if rPage != "" {
		returnPage = rPage
	}

	return c.Render(pkg, ownerList, repos, titleName, returnPage, currentUser)
}

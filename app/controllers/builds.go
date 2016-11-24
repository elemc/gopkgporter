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

	var branches []models.RepoType
	dbgorm.Db.Order("rt_name ASC", true).Find(&branches)

	return c.Render(builds, currentUser, branches)
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

func (c Builds) Push(id uint, branch string) revel.Result {
	if branch == "" {
		branch = "updates"
	}
	currentUser := connected(c.RenderArgs, c.Session)
	if currentUser == nil {
		c.Flash.Error(dontPerm)
		return c.Redirect("/builds")
	}
	var build models.BuildedPackage

	revel.INFO.Printf("Build ID: %d", id)
	revel.INFO.Printf("Branch name: %s", branch)
	ctx := dbgorm.Db.First(&build, "build_id=?", id)
	if ctx.Error != nil {
		c.Flash.Error("Error in fetch build by id: %s", ctx.Error)
		return c.Redirect("/builds")
	}

	dbgorm.Db.Model(&build).Related(&build.BuildPackage, "BuildPackage")
	dbgorm.Db.Model(&build.BuildPackage).Related(&build.BuildPackage.PkgRepo, "PkgRepo")
	dbgorm.Db.Model(&build).Related(&build.Owner, "Owner")
	dbgorm.Db.Model(&build).Related(&build.User, "PushUser")
	dbgorm.Db.Model(&build).Related(&build.PushRepoType, "PushRepoType")

	if (build.Owner.OwnerName != currentUser.UserName) && (currentUser.UserGroup < models.GroupPusher) {
		c.Flash.Error(dontPerm)
		return c.Redirect("/builds")
	}

	var branchRT models.RepoType
	ctx = dbgorm.Db.First(&branchRT, "rt_name=?", branch)
	if ctx.Error != nil {
		revel.ERROR.Printf("Error fetch branch by name: %s", ctx.Error)
		c.Flash.Error(ctx.Error.Error())
		return c.Redirect("/builds")
	}

	push := models.PackagesToPush{}
	push.Fill(build, branch)
	ctx = dbgorm.Db.Create(&push)
	if ctx.Error != nil {
		revel.ERROR.Printf("Error in create push: %s", ctx.Error)
		c.Flash.Error(ctx.Error.Error())
		return c.Redirect("/builds")
	}

	build.Pushed = true
	build.PushUser = *currentUser
	build.PushUserID = currentUser.ID
	build.PushRepoType = branchRT
	build.PushRepoTypeID = branchRT.ID
	build.User = *currentUser
	ctx = dbgorm.Db.Save(&build)
	if ctx.Error != nil {
		revel.ERROR.Printf("Error in save build package: %s", ctx.Error)
		c.Flash.Error(ctx.Error.Error())
		return c.Redirect("/builds")
	}

	return c.Redirect("/builds")
}

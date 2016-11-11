package controllers

import (
	"fmt"
	"gopkgporter/app/models"

	"github.com/revel/revel"
)

// Repos controller
type Repos struct {
	*revel.Controller
}

// Index returns list repositories
func (c Repos) Index() revel.Result {
	var repos []models.Repo
	q := dbgorm.Db.Order("repo_name ASC", true).Find(&repos)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	return c.Render(repos)
}

// Delete function soft delete repository with specified id
func (c Repos) Delete(id int) revel.Result {
	var repo models.Repo
	q := dbgorm.Db.Find(&repo, id)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}
	q = dbgorm.Db.Delete(&repo)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	c.Flash.Success("Repository with new name %s successfully deleted.", repo.RepoName)
	return c.Redirect("/repos")
}

// Insert function create new repository
func (c Repos) Insert() revel.Result {
	repoName := c.Params.Get("RepoName")
	if repoName == "" {
		return c.RenderError(fmt.Errorf("Missing parameter \"Repository name\" or it was empty."))
	}

	repo := models.Repo{RepoName: repoName}
	q := dbgorm.Db.Create(&repo)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	c.Flash.Success("Repository created with new name %s", repoName)
	return c.Redirect("/repos")
}

// Edit function save new repository name
func (c Repos) Edit(id int) revel.Result {
	repoName := c.Params.Get("RepoName")
	if repoName == "" {
		return c.RenderError(fmt.Errorf("Missing parameter \"Repository name\" or it was empty."))
	}

	var repo models.Repo
	q := dbgorm.Db.Find(&repo, id)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	repo.RepoName = repoName
	q = dbgorm.Db.Save(&repo)
	if q.Error != nil {
		return c.RenderError(q.Error)
	}

	c.Flash.Success("Repository saved with new name %s", repoName)
	return c.Redirect("/repos")
}

package controllers

import (
	"encoding/base64"
	"gopkgporter/app/common"
	"gopkgporter/app/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/revel/revel"
)

type Db struct {
	Db     *gorm.DB
	Driver string
	Spec   string
}

var (
	dbgorm Db
)

func InitDB() {

	var err error
	dbgorm.Db, err = common.GetGORM()
	if err != nil {
		revel.ERROR.Fatalf("Connection to database error: %s", err)
		return
	}
	//dbgorm.Db.LogMode(true)

	dbgorm.Db.AutoMigrate(&models.RepoType{}, &models.Repo{}, &models.User{},
		&models.Owner{}, &models.Package{}, &models.BuildedPackage{}, &models.Log{})

	initUsers()
	initRepos()
	initRepoTypes()
}

func initUsers() {
	var users []models.User
	dbgorm.Db.Find(&users)
	if len(users) == 0 {
		user := new(models.User)
		user.UserName = "admin"
		user.UserEMail = "admin@example.com"

		bcryptPass, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		user.UserHashPwd = base64.URLEncoding.EncodeToString(bcryptPass)

		dbgorm.Db.Create(user)
	}
}

func initRepos() {
	var repos []models.Repo
	dbgorm.Db.Find(&repos)
	if len(repos) == 0 {
		repo := new(models.Repo)
		repo.ID = 0
		repo.RepoName = "unknown"
		dbgorm.Db.Create(repo)
	}
}

func initRepoTypes() {
	var repoTypes []models.RepoType
	dbgorm.Db.Find(&repoTypes)
	if len(repoTypes) == 0 {
		repoType := new(models.RepoType)
		repoType.RTName = "releases"
		dbgorm.Db.Create(repoType)

		repoType = new(models.RepoType)
		repoType.RTName = "updates-testing"
		dbgorm.Db.Create(repoType)

		repoType = new(models.RepoType)
		repoType.RTName = "updates"
		dbgorm.Db.Create(repoType)
	}
}

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
	dbgorm.Db.LogMode(true)
	dbgorm.Db.Model(&models.Log{}).Related(&models.User{}, "User")

	dbgorm.Db.AutoMigrate(&models.RepoType{}, &models.Repo{}, &models.User{},
		&models.Owner{}, &models.Package{}, &models.BuildedPackage{}, &models.Log{})

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

	// Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}
	// setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
	// 	for col, size := range colSizes {
	// 		t.ColMap(col).MaxSize = size
	// 	}
	// }
	//
	// t := Dbm.AddTable(models.User{}).SetKeys(true, "UserID")
	// setColumnSizes(t, map[string]int{
	// 	"UserName":    50,
	// 	"UserHashPwd": 100,
	// 	"UserEMail":   100,
	// })
	//
	// t = Dbm.AddTable(models.Owner{}).SetKeys(true, "OwnerID")
	// setColumnSizes(t, map[string]int{
	// 	"OwnerName": 50,
	// })
	//
	// t = Dbm.AddTable(models.RepoType{}).SetKeys(true, "RTID")
	// setColumnSizes(t, map[string]int{
	// 	"RTName": 50,
	// })
	//
	// t = Dbm.AddTable(models.Repo{}).SetKeys(true, "RepoID")
	// setColumnSizes(t, map[string]int{
	// 	"RepoName": 50,
	// })
	//
	// t = Dbm.AddTable(models.BuildedPackage{}).SetKeys(true, "BuildID")
	// setColumnSizes(t, map[string]int{
	// 	"Version": 10,
	// 	"Release": 20,
	// 	"Epoch":   10,
	// 	"TagName": 25,
	// })
	//
	// Dbm.TraceOn("[gorp]", r.INFO)
	//
	// err := Dbm.CreateTables()
	// if err != nil {
	// 	log.Printf("Failed create tables")
	// }
}

//
// type GorpController struct {
// 	*r.Controller
// 	Txn *gorp.Transaction
// }
//
// func (c *GorpController) Begin() r.Result {
// 	txn, err := Dbm.Begin()
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.Txn = txn
// 	return nil
// }
//
// func (c *GorpController) Commit() r.Result {
// 	if c.Txn == nil {
// 		return nil
// 	}
// 	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
// 		panic(err)
// 	}
// 	c.Txn = nil
// 	return nil
// }
//
// func (c *GorpController) Rollback() r.Result {
// 	if c.Txn == nil {
// 		return nil
// 	}
// 	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
// 		panic(err)
// 	}
// 	c.Txn = nil
// 	return nil
// }

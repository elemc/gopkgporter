package common

import (
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

// GetGORM function returns pointer to gorm.DB
func GetGORM() (dbgorm *gorm.DB, err error) {
	var (
		found  bool
		driver string
		spec   string
	)
	if driver, found = revel.Config.String("db.driver"); !found {
		revel.ERROR.Printf("No db.driver found.")
		return
	}
	if spec, found = revel.Config.String("db.spec"); !found {
		revel.ERROR.Printf("No db.spec found.")
		return
	}
	dbgorm, err = gorm.Open(driver, spec)
	if err != nil {
		revel.ERROR.Printf("Failed database open: %s", err)
		return
	}
	dbgorm.DB().SetMaxOpenConns(5)
	dbgorm.DB().SetMaxIdleConns(-1)

	return
}

package main

import (
	"gopkgporter/app/models"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Please use it %s username password", os.Args[0])
	}

	username := os.Args[1]
	password := os.Args[2]

	db, err := GetGORM()
	if err != nil {
		log.Fatalf("Not connected to database")
	}

	user := models.User{}
	user.UserName = username
	hash := user.GeneratePasswordHash(password)
	user.SetPasswordHash(hash)

	ctx := db.Create(&user)
	if ctx.Error != nil {
		log.Fatalf("Error in creating user: %s", ctx.Error)
	}
	os.Exit(0)
}

// GetGORM function returns pointer to gorm.DB
func GetGORM() (dbgorm *gorm.DB, err error) {
	driver := "postgres"
	spec := "user=alex dbname=gopp sslmode=disable"
	dbgorm, err = gorm.Open(driver, spec)
	if err != nil {
		revel.ERROR.Printf("Failed database open: %s", err)
		return
	}

	return
}

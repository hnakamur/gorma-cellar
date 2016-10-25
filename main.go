package main

import (
	"fmt"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/gorma-cellar/app"
	"github.com/goadesign/gorma-cellar/models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/inconshreveable/log15.v2"
)

var db *gorm.DB
var logger log15.Logger
var adb *models.AccountDB
var bdb *models.BottleDB

func main() {

	// Create service
	var err error
	url := fmt.Sprintf("dbname=gorma user=gorma password=gorma sslmode=disable port=%d host=%s", 5432, "localhost")
	fmt.Println(url)
	db, err = gorm.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	db.DropTable(&models.Account{}, &models.Bottle{})
	db.AutoMigrate(&models.Account{}, &models.Bottle{})

	adb = models.NewAccountDB(db)
	bdb = models.NewBottleDB(db)
	db.DB().SetMaxOpenConns(50)
	// Create service
	service := goa.New("API")

	// Setup middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "account" controller
	c := NewAccountController(service)
	app.MountAccountController(service, c)
	// Mount "bottle" controller
	c2 := NewBottleController(service)
	app.MountBottleController(service, c2)
	// Start service, listen on port 8080
	service.ListenAndServe(":8080")
}

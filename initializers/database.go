package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
    var err error

    dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatal("Failed to connect to database")
    }

	// AutoMigrate creates the table if it doesn't exist
	err = DB.AutoMigrate(
		&model.User{},
		&model.Routes{},
		&model.Billing{},
		&model.Purchases{},
		&model.Inventory{},
		&model.Expense{},
	)
	if err != nil {
    	log.Fatal("Failed to migrate database: ", err)
	}

	insertDefaultRoutes()
}

func insertDefaultRoutes() {
	var count int64
	DB.Model(&model.Routes{}).Count(&count)
	if count > 0 {
		return // routes already exist, skip insertion
	}

	defaultRoutes := []model.Routes{
		{Route: "/billing", Heading: "Billing", Role: constants.RoleAdmin},
		{Route: "/expenses", Heading: "Expenses", Role: constants.RoleAdmin},
		{Route: "/inventory", Heading: "Inventory", Role: constants.RoleAdmin},
		{Route: "/", Heading: "Logout", Role: constants.RoleAdmin},
		{Route: "/billing", Heading: "Billing", Role: constants.RoleUser},
		{Route: "/", Heading: "Logout", Role: constants.RoleUser},
	}

	if err := DB.Create(&defaultRoutes).Error; err != nil {
		log.Println("Failed to insert default routes:", err)
	} else {
		log.Println("Default routes inserted successfully")
	}
}
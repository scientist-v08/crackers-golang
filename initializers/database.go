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
		&model.PriceList{},
	)
	if err != nil {
    	log.Fatal("Failed to migrate database: ", err)
	}

	insertDefaultRoutes()
	insertStandardPriceList()
}

func insertDefaultRoutes() {
	var count int64
	DB.Model(&model.Routes{}).Count(&count)
	if count == 6 {
		if err := DB.Delete(&model.Routes{}, []int{1,2,3,4,5,6}).Error; err != nil {
			log.Println("Failed to delete previous routes:", err)
		} else {
			log.Println("Deleted Old Routes")
		}
		defaultRoutes := []model.Routes{
			{Route: "/billing", Heading: "Billing", Role: constants.RoleAdmin},
			{Route: "/price-list", Heading: "Price List", Role: constants.RoleAdmin},
			{Route: "/expenses", Heading: "Expenses", Role: constants.RoleAdmin},
			{Route: "/inventory", Heading: "Inventory", Role: constants.RoleAdmin},
			{Route: "/", Heading: "Logout", Role: constants.RoleAdmin},
			{Route: "/billing", Heading: "Billing", Role: constants.RoleUser},
			{Route: "/price-list", Heading: "Price List", Role: constants.RoleUser},
			{Route: "/", Heading: "Logout", Role: constants.RoleUser},
		}

		if err := DB.Create(&defaultRoutes).Error; err != nil {
			log.Println("Failed to insert default routes:", err)
		} else {
			log.Println("Default routes inserted successfully")
		} 
	} else {
		if count == 0 {
			defaultRoutes := []model.Routes{
				{Route: "/billing", Heading: "Billing", Role: constants.RoleAdmin},
				{Route: "/price-list", Heading: "Price List", Role: constants.RoleAdmin},
				{Route: "/expenses", Heading: "Expenses", Role: constants.RoleAdmin},
				{Route: "/inventory", Heading: "Inventory", Role: constants.RoleAdmin},
				{Route: "/", Heading: "Logout", Role: constants.RoleAdmin},
				{Route: "/billing", Heading: "Billing", Role: constants.RoleUser},
				{Route: "/price-list", Heading: "Price List", Role: constants.RoleUser},
				{Route: "/", Heading: "Logout", Role: constants.RoleUser},
			}

			if err := DB.Create(&defaultRoutes).Error; err != nil {
				log.Println("Failed to insert default routes:", err)
			} else {
				log.Println("Default routes inserted successfully")
			} 
		} else {
			return // routes already exist, skip insertion
		}
	}
}

func insertStandardPriceList() {
	var count int64
	DB.Model(&model.PriceList{}).Count(&count)
	if count == 0 {
		standardPriceList := []model.PriceList{
			{Brand: "Standard - Fancy", Item: "CCB", Price: 350},
			{Brand: "Standard - Fancy", Item: "Golden Whistle Small", Price: 260},
			{Brand: "Standard - Fancy", Item: "Golden Whistle Big", Price: 540},
			{Brand: "Standard - Fancy", Item: "Gold Rush", Price: 510},
			{Brand: "Standard - Fancy", Item: "Aerial Outs", Price: 516},
			{Brand: "Standard - Fancy", Item: "7 Shots", Price: 307},
			{Brand: "Standard - Fancy", Item: "Golden Drops", Price: 620},
			{Brand: "Standard - Fancy", Item: "Magic Whip", Price: 155},
			{Brand: "Standard - Fancy", Item: "Stones", Price: 155},
			{Brand: "Standard - Fancy", Item: "Olympic Torch", Price: 550},
			{Brand: "Standard", Item: "Jet Mix - 12 Shots (8 variants)", Price: 440},
			{Brand: "Standard", Item: "Fly All - 12 Shots (5 variants)", Price: 270},
			{Brand: "Standard", Item: "Touch Sky - 25 Shots (3 variants)", Price: 1000},
			{Brand: "Standard", Item: "Grand Finale - 40 Shots (5 variants)", Price: 1360},
			{Brand: "Standard", Item: "Fly Boss - 60 Shots (4 variants)", Price: 1800},
			{Brand: "Standard", Item: "80's HITS - 80 Shots (3 variants)", Price: 2400},
			{Brand: "Standard", Item: "Sensation - 100 Shots", Price: 4400},
			{Brand: "Standard", Item: "Sun Flower - 100 Shots", Price: 4400},
			{Brand: "Standard", Item: "Moon Light - 100 Shots (4 variants)", Price: 2600},
			{Brand: "Standard", Item: "Crushers - 125 Shots", Price: 5850},
			{Brand: "Standard", Item: "Rainbow Dance - 240 Shots", Price: 6000},
			{Brand: "Standard", Item: "1000 Wala", Price: 930},
			{Brand: "Standard", Item: "2000 Wala", Price: 1750},
			{Brand: "Standard", Item: "5000 Wala", Price: 4100},
			{Brand: "Standard", Item: "10000 Wala", Price: 8300},
			{Brand: "Standard", Item: "1\" Comet - Mix Masala (10 variants)", Price: 100},
			{Brand: "Standard", Item: "1.5\" Comet - Friendly Fire (12 variants)", Price: 190},
			{Brand: "Standard", Item: "1.75\" Comet - Team Captain (12 variants)", Price: 230},
			{Brand: "Standard", Item: "1.75\" Comet - I Comet", Price: 690},
			{Brand: "Standard", Item: "1.75\" Comet - Space Race (12 variants)", Price: 690},
			{Brand: "Standard", Item: "1.75\" Comet - Thunder Bolt & Jasmine Drops", Price: 750},
			{Brand: "Standard", Item: "2\" Asmaan Guru", Price: 290},
			{Brand: "Standard", Item: "2.5\" Peanut (Double Shot)", Price: 540},
			{Brand: "Standard", Item: "3\" Moon Modi & Kya Baath Hai (13 variants)", Price: 535},
			{Brand: "Standard", Item: "3\" Twitter Glitter (7 variants)", Price: 765},
		}
	
		if err := DB.Create(&standardPriceList).Error; err != nil {
			log.Println("Failed to insert standard price list")
		} else {
			log.Println("Standard Price list inserted")
		}
	}
}
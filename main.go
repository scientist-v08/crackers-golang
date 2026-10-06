package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/routes"
)

var isProduction bool

func setFiberMode() {
	// Set Fiber MODE
	if mode := os.Getenv("GIN_MODE"); mode == "debug" {
		isProduction = false
	} else {
		isProduction = true
	}
}

func init() {
	initializers.LoadEnvVariables()
	setFiberMode()
	initializers.ConnectToDb()
}

func main() {
	r := fiber.New(fiber.Config{
		ServerHeader:  "Vinayaka Crackers",
		AppName:       "Crackers API",
		StrictRouting: true,
		CaseSensitive: true,
	})
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{allowedOrigin},
		AllowMethods:  []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
    	ExposeHeaders: []string{"Content-Disposition", "Content-Length"},
		AllowCredentials: true,
		MaxAge: int(12 * time.Hour),
	}))

	// Now initialize all the routes
	routes.RegisterUserRoutes(r)
	routes.RegisterBillingRoutes(r)
	routes.RegisterInventoryRoutes(r)
	routes.RegisterExpenseRoutes(r)
	routes.RegisterPriceList(r)
	
	// Now run the application
	r.Listen(":" + initializers.Port, fiber.ListenConfig{EnablePrefork: isProduction})
}

package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{allowedOrigin},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
    	ExposeHeaders: []string{"Content-Disposition", "Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	// Now initialize all the routes
	routes.RegisterUserRoutes(r)
	routes.RegisterBillingRoutes(r)
	routes.RegisterInventoryRoutes(r)
	routes.RegisterExpenseRoutes(r)
	
	// Now run the application
	r.Run()
}

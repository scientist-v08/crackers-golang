package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	if os.Getenv("APP_ENV") == "production" {
		log.Println("Production mode: using system environment variables")
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}
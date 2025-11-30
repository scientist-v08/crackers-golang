package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	env := os.Getenv("APP_ENV")
	file := ".env"

	switch env {
		case "production":
			file = ".env.production"
		case "staging":
			file = ".env.staging"
	}
	err := godotenv.Load(file)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
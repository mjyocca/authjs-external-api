package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	envFileName := ".env.local"
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

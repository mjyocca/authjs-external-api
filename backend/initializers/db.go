package initializers

import (
	"fmt"
	"os"

	"github.com/mjyocca/authjs-external-api/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func getPostgresUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func ConnectToDatabase() {
	var err error
	dbURL := getPostgresUrl()
	fmt.Println(dbURL)
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
}

func AutoMigrate() {
	DB.AutoMigrate(&models.User{})
}

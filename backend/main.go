package main

import (
	"fmt"

	"github.com/mjyocca/authjs-external-api/backend/controllers"
	"github.com/mjyocca/authjs-external-api/backend/initializers"
	"github.com/mjyocca/authjs-external-api/backend/router"
	store "github.com/mjyocca/authjs-external-api/backend/stores"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.AutoMigrate()
}

func main() {

	app := router.New()

	// create stores
	userStore := store.NewUserStore(initializers.DB)

	// create handler
	handler := controllers.NewHandler(userStore)

	// register routes with handlers
	handler.Register(app)

	if err := app.Listen(":8000"); err != nil {
		fmt.Printf("%v", err)
	}
}

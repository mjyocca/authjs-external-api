package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/controllers"
	"github.com/mjyocca/authjs-external-api/backend/initializers"
	"github.com/mjyocca/authjs-external-api/backend/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.RunDBMigrations()
}

func main() {

	app := fiber.New()

	/* jwt auth middleware */
	app.Use(middleware.Auth(middleware.AuthConfig{}))

	api := app.Group("/api", middleware.Adapter)

	/* adapter routes */
	api.Route("/adapter", func(route fiber.Router) {
		route.Get("/user", controllers.GetUser)
		route.Post("/user", controllers.CreateUser)
		route.Patch("/user", controllers.LinkAccount)
	})
	api.Get("/user", controllers.UserIndex)

	app.Listen(":8000")
}

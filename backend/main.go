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

	/* public routes */
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	/* protected ~ jwt auth middleware */
	api := app.Group("/api", middleware.Auth(middleware.AuthConfig{}))
	nextAuthAdapter := api.Group("/adapter", middleware.Adapter)

	/* NextAuth Adapter routes */
	nextAuthAdapter.Route("/", func(route fiber.Router) {
		route.Get("/user", controllers.GetUserAdapter)
		route.Post("/user", controllers.CreateUserAdapter)
		route.Patch("/user", controllers.LinkAccountAdapter)
	})

	/* API routes */
	api.Route("/", func(route fiber.Router) {
		api.Get("/user", controllers.UserIndex)
	})

	app.Use(middleware.NotFound)

	app.Listen(":8000")
}

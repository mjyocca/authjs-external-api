package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/controllers"
	"github.com/mjyocca/authjs-external-api/backend/initializers"
	"github.com/mjyocca/authjs-external-api/backend/middleware"
	store "github.com/mjyocca/authjs-external-api/backend/stores"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.AutoMigrate()
}

func main() {

	app := fiber.New()

	/* create stores */
	userStore := store.NewUserStore(initializers.DB)

	/* init handler */
	handler := controllers.NewHandler(userStore)

	/* public routes */
	app.Get("/health", handler.HealthCheck)

	/* protected ~ jwt auth middleware */
	api := app.Group("/api", middleware.Auth(middleware.AuthConfig{}))
	nextAuthAdapter := api.Group("/adapter", middleware.Adapter)

	/* NextAuth Adapter routes */
	nextAuthAdapter.Route("/", func(route fiber.Router) {
		route.Get("/user", handler.GetUserAdapter)
		route.Post("/user", handler.CreateUserAdapter)
		route.Patch("/user", handler.LinkAccountAdapter)
	})

	/* API routes */
	api.Route("/", func(route fiber.Router) {
		api.Get("/user", handler.CurrentUser)
	})

	app.Use(controllers.NotFound)

	app.Listen(":8000")
}

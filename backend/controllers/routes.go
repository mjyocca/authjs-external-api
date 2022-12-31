package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/middleware"
)

func (h *Handler) Register(app *fiber.App) {

	/* public routes */
	app.Get("/health", h.HealthCheck)

	/* protected ~ jwt auth middleware */
	api := app.Group("/api", middleware.Auth(middleware.AuthConfig{}))
	nextAuthAdapter := api.Group("/adapter", middleware.Adapter)

	/* NextAuth Adapter routes */
	nextAuthAdapter.Route("/", func(route fiber.Router) {
		route.Get("/user", h.GetUserAdapter)
		route.Post("/user", h.CreateUserAdapter)
		route.Patch("/user", h.LinkAccountAdapter)
	})

	/* API routes */
	api.Route("/", func(route fiber.Router) {
		api.Get("/user", h.CurrentUser)
	})

	/* Catch all */
	app.Use(h.NotFound)

}

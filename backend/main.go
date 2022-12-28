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

	/* routes */
	app.Get("/", func(c *fiber.Ctx) error {
		// claimData := c.Locals("jwtClaims")
		// fmt.Println(claimData)
		return c.JSON(fiber.Map{
			"data": "Hello, World!",
		})
	})

	app.Get("/user", controllers.UserIndex)

	app.Listen(":8000")
}

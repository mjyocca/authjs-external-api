package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/initializers"
	auth "github.com/mjyocca/authjs-external-api/backend/internal/auth"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {

	app := fiber.New()

	/* jwt middleware */
	app.Use(
		auth.NewConfig(auth.Config{}),
	)

	/* routes */
	app.Get("/", func(c *fiber.Ctx) error {
		// claimData := c.Locals("jwtClaims")
		// fmt.Println(claimData)
		return c.JSON(fiber.Map{
			"data": "Hello, World!",
		})
	})

	app.Listen(":8000")
}

package main

import (
	fiber "github.com/gofiber/fiber/v2"
	jwechecker "github.com/mjyocca/authjs-external-api/jwechecker"
)

func main() {
	app := fiber.New()

	/* jwt middleware */
	app.Use(
		jwechecker.NewConfig(jwechecker.Config{}),
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

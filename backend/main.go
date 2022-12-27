package main

import (
	fiber "github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	/* jwt middleware */
	app.Use(
		NewJweConfig(JweConfig{}),
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

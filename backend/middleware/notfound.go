package middleware

import fiber "github.com/gofiber/fiber/v2"

func NotFound(c *fiber.Ctx) error {
	return c.SendStatus(404)
}

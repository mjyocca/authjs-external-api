package middleware

import (
	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
)

func Adapter(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	if userClaims["Role"] != "adapter" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}

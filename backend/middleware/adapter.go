package middleware

import (
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
)

func Adapter(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	if userClaims["Role"] != "adapter" {
		fmt.Printf("Unauthorized: Role not adapter")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}

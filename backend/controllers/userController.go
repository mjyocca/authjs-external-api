package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mjyocca/authjs-external-api/backend/initializers"
	"github.com/mjyocca/authjs-external-api/backend/models"
)

func UserIndex(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	providerId := claims["ProviderAccountId"].(string)
	providerType := claims["Provider"].(string)

	user := new(models.User)
	if providerType == "github" {
		initializers.DB.Find(&user, "github_id = ?", providerId)
	}

	return c.JSON(userResponse(user))
}

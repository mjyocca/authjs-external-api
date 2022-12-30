package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mjyocca/authjs-external-api/backend/initializers"
	"github.com/mjyocca/authjs-external-api/backend/models"
)

func UserIndex(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	user := new(models.User)
	initializers.DB.Find(&user, "github_id = ?", userClaims["ProviderAccountId"].(string))
	return c.JSON(fiber.Map{
		"data": user,
	})
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	GithubId string `json:"githubId"`
}

func UserPost(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return err
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	providerAccountId := userClaims["ProviderAccountId"].(string)
	providerType := userClaims["Provider"].(string)

	return c.JSON(fiber.Map{
		"Name":       user.Name,
		"Email":      user.Email,
		"provider":   providerType,
		"providerId": providerAccountId,
	})
}

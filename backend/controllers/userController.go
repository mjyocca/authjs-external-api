package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func UserIndex(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	return c.JSON(fiber.Map{
		"id": userClaims["Id"].(string),
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

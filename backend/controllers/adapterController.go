package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/initializers"
	"github.com/mjyocca/authjs-external-api/backend/models"
)

type createUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
}

func CreateUserAdapter(c *fiber.Ctx) error {
	log.Println("adapter CreateUser")

	u := new(createUser)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	/* If already exists return existing user */
	existingUser := models.User{}
	initializers.DB.Where(&models.User{Email: u.Email}).First(&existingUser)

	if existingUser != (models.User{}) {
		fmt.Println("Existing User: ", existingUser)
		return c.JSON(fiber.Map{
			"id":    existingUser.ID,
			"email": existingUser.Email,
		})
	}

	log.Println("Creating new user")
	user := models.User{Name: u.Name, Email: u.Email, Image: u.Image}
	initializers.DB.Create(&user)

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	})
}

/* options: byId, byEmail, byAccount */
func GetUserAdapter(c *fiber.Ctx) error {
	log.Println("adapter GetUser")

	userId := c.Query("userId")
	email := c.Query("email")
	account := c.Query("accountId")

	user := models.User{}
	if userId != "" {
		initializers.DB.First(&user, userId)
		if user == (models.User{}) {
			return c.JSON(fiber.Map{"msg": "Not Found"})
		}
		return c.JSON(user)
	}

	if email != "" {
		initializers.DB.Where(&models.User{Email: email}).First(&user)
		if user == (models.User{}) {
			return c.JSON(fiber.Map{"msg": "Not Found"})
		}
		return c.JSON(user)
	}

	if account != "" {
		initializers.DB.Where(&models.User{GithubId: account}).First(&user)
		if user == (models.User{}) {
			return c.JSON(fiber.Map{"msg": "Not Found"})
		}
		return c.JSON(fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"image": user.Image,
		})
	}

	return c.JSON(fiber.Map{
		"msg": "Not Found",
	})
}

type linkAccountPayload struct {
	Provider          string
	Type              string
	ProviderAccountId string
	AccessToken       string `json:"access_token"`
	TokenType         string `json:"token_type"`
	Scope             string
	UserId            int64 `json:"userId"`
}

func LinkAccountAdapter(c *fiber.Ctx) error {
	log.Println("adapter LinkAccount")
	link := new(linkAccountPayload)
	if err := c.BodyParser(link); err != nil {
		return err
	}

	user := models.User{}
	initializers.DB.First(&user, link.UserId)

	/* user is not found */
	if user == (models.User{}) {
		return c.JSON(fiber.Map{
			"msg": "User not found",
		})
	}

	/* update current user with provider link */
	user.GithubId = link.ProviderAccountId
	initializers.DB.Save(&user)

	return c.JSON(fiber.Map{
		"id":       user.ID,
		"githubId": user.GithubId,
	})
}

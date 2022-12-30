package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	return c.JSON(user)
}

/* options: byId, byEmail, byAccount */
func GetUserAdapter(c *fiber.Ctx) error {
	log.Println("adapter GetUser")

	userId := c.Query("userId")
	email := c.Query("email")
	account := c.Query("accountId")

	user := models.User{}
	if userId != "" {
		externalId, err := uuid.Parse(userId)
		if err != nil {
			return c.JSON(errorResponse("server error"))
		}
		initializers.DB.Where(&models.User{ExternalID: externalId}).First(&user)
		if user == (models.User{}) {
			return c.JSON(notFoundResponse())
		}
		return c.JSON(userResponse(&user))
	}

	if email != "" {
		initializers.DB.Where(&models.User{Email: email}).First(&user)
		if user == (models.User{}) {
			return c.JSON(notFoundResponse())
		}
		return c.JSON(userResponse(&user))
	}

	if account != "" {
		// to-do: need to add checks for provider type
		initializers.DB.Where(&models.User{GithubId: account}).First(&user)
		if user == (models.User{}) {
			return c.JSON(notFoundResponse())
		}
		return c.JSON(userResponse(&user))
	}

	return c.JSON(notFoundResponse())
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
			"msg": "Not Found",
		})
	}

	/* update current user with provider link */
	user.GithubId = link.ProviderAccountId
	initializers.DB.Save(&user)

	return c.JSON(userResponse(&user))
}

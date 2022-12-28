package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	log.Println("adapter CreateUser")
	return c.JSON(fiber.Map{})
}

/* options: byId, byEmail, byAccount */
func GetUser(c *fiber.Ctx) error {
	log.Println("adapter GetUser")
	return c.JSON(fiber.Map{})
}

func LinkAccount(c *fiber.Ctx) error {
	log.Println("adapter LinkAccount")
	return c.JSON(fiber.Map{})
}

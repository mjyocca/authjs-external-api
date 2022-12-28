package controllers

import "github.com/gofiber/fiber/v2"

func UserIndex(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"id": "001",
	})
}

package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/models"
)

func userResponse(user *models.User) fiber.Map {
	return fiber.Map{
		"data":   user,
		"meta":   fiber.Map{},
		"status": "success",
	}
}

func errorResponse(msg string) fiber.Map {
	return fiber.Map{
		"msg":    msg,
		"status": "error",
	}
}

func notFoundResponse() fiber.Map {
	return fiber.Map{
		"msg":    "Not Found",
		"status": "warning",
	}
}

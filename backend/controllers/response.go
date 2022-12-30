package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/models"
)

func userResponse(user *models.User) fiber.Map {
	return fiber.Map{
		"data": user,
		"stats": fiber.Map{
			"hi": "there",
		},
	}
}

func errorResponse(msg string) fiber.Map {
	return fiber.Map{
		"msg": msg,
	}
}

func notFoundResponse() fiber.Map {
	return fiber.Map{
		"msg": "Not Found",
	}
}

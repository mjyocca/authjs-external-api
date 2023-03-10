package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/models"
)

type ResponseMessage string

const (
	NotFound            ResponseMessage = "Not Found"
	InternalServerError ResponseMessage = "Internal Server Error"
	CannotProcessEntity ResponseMessage = "Cannot Process Request"
	AccessForbidden     ResponseMessage = "Access is Forbidden"
)

func userResponse(user *models.User) fiber.Map {
	return fiber.Map{
		"data":   user,
		"meta":   fiber.Map{},
		"status": "success",
	}
}

func errorResponse(msg ResponseMessage) fiber.Map {
	return fiber.Map{
		"msg":    msg,
		"status": "error",
	}
}

func notFoundResponse() fiber.Map {
	return fiber.Map{
		"msg":    NotFound,
		"status": "warning",
	}
}

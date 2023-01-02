package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mjyocca/authjs-external-api/backend/helpers"
)

func (h *Handler) CurrentUser(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)

	userId, err := helpers.GetUUID(claims["Id"].(string))
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse("cannot process request"))
	}

	user, err := h.userStore.GetByExternalID(userId)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse("cannot process request"))
	}
	if user == nil {
		return c.Status(http.StatusForbidden).JSON(errorResponse("access is forbidden"))
	}
	return c.Status(http.StatusFound).JSON(userResponse(user))
}

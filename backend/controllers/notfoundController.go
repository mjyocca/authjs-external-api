package controllers

import fiber "github.com/gofiber/fiber/v2"

func (h *Handler) NotFound(c *fiber.Ctx) error {
	return c.SendStatus(404)
}

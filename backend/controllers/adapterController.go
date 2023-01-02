package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mjyocca/authjs-external-api/backend/helpers"
	model "github.com/mjyocca/authjs-external-api/backend/models"
)

// Next-Auth/Authjs Adapter function
func (h *Handler) CreateUserAdapter(c *fiber.Ctx) error {
	var u model.User
	req := &userCreateAdapterRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse(CannotProcessEntity))
	}
	existingUser, _ := h.userStore.GetByEmail(req.Email)
	if *existingUser != (model.User{}) {
		return c.Status(http.StatusAlreadyReported).JSON(userResponse(existingUser))
	}
	// populate data
	u.Name = req.Name
	u.Email = req.Email
	u.Image = req.Image

	if err := h.userStore.Create(&u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse(CannotProcessEntity))
	}
	return c.Status(http.StatusCreated).JSON(userResponse(&u))
}

// Next-Auth/Authjs Adapter function
func (h *Handler) GetUserAdapter(c *fiber.Ctx) error {
	// query params
	userId := c.Query("id")
	email := c.Query("email")
	providerId := c.Query("providerId")
	providerType := c.Query("providerType")

	// or conditions to search by
	conditions := map[string]interface{}{}
	if userId != "" {
		conditions["external_id"] = userId
	}
	if email != "" {
		conditions["email"] = email
	}
	if providerId != "" && providerType != "" {
		fieldName := helpers.ProviderFieldMapping[providerType]
		conditions[fieldName] = providerId
	}
	u, err := h.userStore.GetUserByORConditions(conditions)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse(InternalServerError))
	}
	if u == nil {
		return c.Status(http.StatusForbidden).JSON(errorResponse(AccessForbidden))
	}

	return c.Status(http.StatusFound).JSON(userResponse(u))
}

// Next-Auth/Authjs Adapter function
func (h *Handler) LinkAccountAdapter(c *fiber.Ctx) error {
	req := &linkAccountAdapterRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse(CannotProcessEntity))
	}

	// parse user id to uuid format
	userId, err := helpers.GetUUID(req.UserId)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse(CannotProcessEntity))
	}

	// get existing user from database
	user, err := h.userStore.GetByExternalID(userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse(InternalServerError))
	}
	if user == nil {
		return c.Status(http.StatusNotFound).JSON(notFoundResponse())
	}

	// update user field(s)
	switch provider := req.Provider; provider {
	case string(helpers.Github):
		user.GithubId = req.ProviderAccountId
	case string(helpers.Google):
		user.GoogleId = req.ProviderAccountId
	}

	if err := h.userStore.Update(user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse(CannotProcessEntity))
	}

	return c.Status(http.StatusOK).JSON(userResponse(user))
}

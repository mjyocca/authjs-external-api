package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	model "github.com/mjyocca/authjs-external-api/backend/models"
)

// Next-Auth/Authjs Adapter function
func (h *Handler) CreateUserAdapter(c *fiber.Ctx) error {
	var u model.User
	req := &userCreateRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse("cannot process request"))
	}
	existingUser, _ := h.userStore.GetByEmail(req.User.Email)
	if existingUser != (&model.User{}) {
		return c.Status(http.StatusAlreadyReported).JSON(&existingUser)
	}
	if err := h.userStore.Create(&u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse("cannot process request"))
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
		conditions["ExternalID"] = userId
	}
	if email != "" {
		conditions["Email"] = email
	}
	if providerId != "" && providerType != "" {
		conditions[providerType] = providerId
	}
	u, err := h.userStore.GetUserByORConditions(conditions)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse("internal server error"))
	}
	if u == nil {
		return c.Status(http.StatusForbidden).JSON(errorResponse("access is forbidden"))
	}

	return c.Status(http.StatusFound).JSON(userResponse(u))
}

// Next-Auth/Authjs Adapter function
func (h *Handler) LinkAccountAdapter(c *fiber.Ctx) error {
	req := &userLinkAccountRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse("cannot process request"))
	}

	// parse user id to uuid format
	userId, err := getUUID(req.UserId)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse("cannot process request"))
	}

	// get existing user from database
	user, err := h.userStore.GetByExternalID(userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse("internal server error"))
	}
	if user == nil {
		return c.Status(http.StatusNotFound).JSON(notFoundResponse())
	}

	// update user field(s)
	switch provider := req.Provider; provider {
	case "Github":
		user.GithubId = req.ProviderAccountId
		// case "Google":
	}

	if err := h.userStore.Update(user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(errorResponse("cannot process request"))
	}

	return c.Status(http.StatusOK).JSON(userResponse(user))
}

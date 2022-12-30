package controllers

import (
	store "github.com/mjyocca/authjs-external-api/backend/stores"
)

type Handler struct {
	userStore store.User
}

func NewHandler(us store.User) *Handler {
	return &Handler{
		userStore: us,
	}
}

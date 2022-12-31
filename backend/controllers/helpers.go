package controllers

import "github.com/google/uuid"

// Converts stringified UUID to proper UUID datatype
func getUUID(userId string) (uuid.UUID, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.New(), err
	}
	return id, nil
}

// move to utils/helpers package
var providerFieldMapping = map[string]string{
	"github": "github_id",
	"google": "google_id",
}

type NextAuthProvider string

const (
	Github NextAuthProvider = "github"
	Google NextAuthProvider = "google"
)

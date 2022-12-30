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

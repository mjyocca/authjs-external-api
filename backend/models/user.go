package models

import (
	"encoding/json"
	"time"
)

type User struct {
	// gorm.Model
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt"`

	Name     string `json:"name"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	GithubId string `json:"githubId"`
}

func (u *User) Providers() []string {
	providers := []string{}
	if u.GithubId != "" {
		providers = append(providers, "github")
	}
	return providers
}

func (u *User) MarshalJSON() ([]byte, error) {
	type UserAlias User
	return json.Marshal(&struct {
		*UserAlias
		Providers []string `json:"providers"`
	}{
		UserAlias: (*UserAlias)(u),
		Providers: u.Providers(),
	})
}

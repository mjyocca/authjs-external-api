package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt"`

	ExternalID uuid.UUID `gorm:"type:uuid;index" json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Image      string    `json:"image"`
	GithubId   string    `json:"githubId"`
	GoogleId   string    `json:"googleId"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ExternalID = uuid.New()
	return
}

func (user *User) Providers() []string {
	providers := []string{}
	if user.GithubId != "" {
		providers = append(providers, "github")
	}
	return providers
}

func (user *User) MarshalJSON() ([]byte, error) {
	type UserAlias User
	return json.Marshal(&struct {
		*UserAlias
		Providers []string `json:"providers"`
	}{
		UserAlias: (*UserAlias)(user),
		Providers: user.Providers(),
	})
}

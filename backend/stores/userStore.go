package store

import (
	"errors"

	"github.com/google/uuid"
	model "github.com/mjyocca/authjs-external-api/backend/models"
	"gorm.io/gorm"
)

type User interface {
	GetByExternalID(externalId uuid.UUID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByProviderID(providerId string, providerType string) (*model.User, error)
	GetUserByORConditions(fields map[string]interface{}) (*model.User, error)
	Create(u *model.User) (err error)
	Update(u *model.User) error
}

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetByExternalID(externalId uuid.UUID) (*model.User, error) {
	var m model.User
	if err := us.db.Find(&m, &model.User{ExternalID: externalId}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetByEmail(email string) (*model.User, error) {
	var m model.User
	if err := us.db.Find(&m, &model.User{Email: email}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetByProviderID(providerId string, providerType string) (*model.User, error) {
	var m model.User
	if err := us.db.Find(&m, &model.User{GithubId: providerId}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetUserByORConditions(fields map[string]interface{}) (*model.User, error) {
	var m model.User

	count := 0
	for key, value := range fields {
		clause := map[string]interface{}{
			key: value,
		}
		if count <= 0 {
			us.db.Where(clause)
		} else {
			us.db.Or(clause)
		}
		count++
	}

	if err := us.db.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) Create(u *model.User) (err error) {
	return us.db.Create(u).Error
}

func (us *UserStore) Update(u *model.User) error {
	return us.db.Model(u).Updates(u).Error
}

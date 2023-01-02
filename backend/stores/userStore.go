package store

import (
	"errors"

	"github.com/google/uuid"
	helpers "github.com/mjyocca/authjs-external-api/backend/helpers"
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
	condition := &model.User{}

	if providerType == string(helpers.Github) {
		condition.GithubId = providerId
	} else if providerType == string(helpers.Google) {
		condition.GoogleId = providerId
	}

	if err := us.db.Find(&m, condition).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *UserStore) GetUserByORConditions(fields map[string]interface{}) (*model.User, error) {
	var m model.User
	queryDB := us.db.Session(&gorm.Session{})

	count := 0
	for key, value := range fields {
		clause := map[string]interface{}{
			key: value,
		}
		if count <= 0 {
			queryDB.Where(clause)
		} else {
			queryDB.Or(clause)
		}
		count++
	}

	// think this is the problem
	if err := queryDB.First(&m).Error; err != nil {
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

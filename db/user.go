package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"leggett.dev/devmarks/api/model"
)

// GetUserByEmail returns the user with the specified email address from the database.
func (db *Database) GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	if err := db.First(&user, model.User{Email: email}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "unable to get user")
	}
	return &user, nil
}

// CreateUser inserts a new user into the database.
func (db *Database) CreateUser(user *model.User) error {
	return db.Create(user).Error
}

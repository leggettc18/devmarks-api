package model

import (
	"golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash takes a plaintext password and generates a bcryt
// hash of it.
func GeneratePasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// ComparePasswordHash takes a password hash and a plaintext password and returns true
// if the plaintext password hashes into the password hash.
func ComparePasswordHash(hashedPassword, givenPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, givenPassword)
	return err == nil
}

// User is a model representing the Users our app can save. They can own any number of bookmarks,
// any number of folders, and be a part of any number of organizations.
type User struct {
	Model

	Email          string `json:"email"`
	HashedPassword []byte `json:"-"`

	Bookmarks []Bookmark `gorm:"foreignkey:OwnerID"`
}

// SetPassword takes a plaintext password and saves the resulting hash to the User model.
func (u *User) SetPassword(password string) error {
	hashed, err := GeneratePasswordHash([]byte(password))
	if err != nil {
		return err
	}
	u.HashedPassword = hashed
	return nil
}

// CheckPassword takes a plaintext password and compares it to the User's hashed password.
func (u *User) CheckPassword(password string) bool {
	return ComparePasswordHash(u.HashedPassword, []byte(password))
}

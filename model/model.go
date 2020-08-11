package model

import (
	"crypto/rand"
	"time"
)

// ID is a wrapper around []byte representing the ID of any of our requests.
type ID []byte

// NewID generates a new random ID value.
func NewID() ID {
	ret := make(ID, 20)
	if _, err := rand.Read(ret); err != nil {
		panic(err)
	}
	return ret
}

// Model is a base model to provide certain values to all of our other models.
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	Name           string
	Email          string
	HashedPassword []byte

	Bookmarks []Bookmark
}

package model

import "github.com/jinzhu/gorm"

type Tag struct {
	gorm.Model

	Name string

	Bookmarks []Bookmark `gorm:"many2many:bookmark_tag;"`
}

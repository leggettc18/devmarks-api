package model

import "github.com/jinzhu/gorm"

type Bookmark struct {
	gorm.Model

	Name  string
	Url   string
	Color string

	User    User
	UserID  uint
	Folders []Folder `gorm:"many2many:bookmark_folder;"`
	Tags    []Tag    `gorm:"many2many:bookmark_tag;"`
}

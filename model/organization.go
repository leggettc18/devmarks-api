package model

import "github.com/jinzhu/gorm"

type Organization struct {
	gorm.Model

	Name        string
	Description string

	Folders []Folder `gorm:"many2many:folder_organization;"`
	Users   []User   `gorm:"many2many:organization_user;"`
	Owner   User     `gorm:"association_foreignkey:UserID"`
}

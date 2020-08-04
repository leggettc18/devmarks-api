package model

type Organization struct {
	Model

	Name        string
	Description string

	Folders []Folder `gorm:"many2many:folder_organization;"`
	Users   []User   `gorm:"many2many:organization_user;"`
	Owner   User     `gorm:"association_foreignkey:UserID"`
}

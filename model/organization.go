package model

// Organization is a model representing the organizations our app can save. Organizations provide
// a way to group any number of users together, can have access to any number of folders,
// and are owned by a single user.
type Organization struct {
	Model

	Name        string
	Description string

	Folders []Folder `gorm:"many2many:folder_organization;"`
	Users   []User   `gorm:"many2many:organization_user;"`
	Owner   User     `gorm:"association_foreignkey:UserID"`
}

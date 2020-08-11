package model

// Folder is a model that represents folders of bookmarks that can exist in our app,
// containing bookmarks, owned by one user, and accessed by individual users and or
// users belonging to a certain organization.
type Folder struct {
	Model

	Name  string
	Color string

	Owner         User           `gorm:"association_foreignkey:UserID"`
	Organizations []Organization `gorm:"many2many:folder_organization;"`
	Bookmarks     []Bookmark     `gorm:"many2many:bookmark_folder;"`
	Users         []User         `gorm:"many2many:folder_user;"`
}

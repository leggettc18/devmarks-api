package model

type Folder struct {
	Model

	Name  string
	Color string

	Owner         User           `gorm:"association_foreignkey:UserID"`
	Organizations []Organization `gorm:"many2many:folder_organization;"`
	Bookmarks     []Bookmark     `gorm:"many2many:bookmark_folder;"`
	Users         []User         `gorm:many2many:folder_user;"`
}

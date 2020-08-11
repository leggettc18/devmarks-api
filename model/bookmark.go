package model

// Bookmark is a model that represents the bookmarks our app can save. They are owned by one user,
// Can be in any number of folders, and can have any number of tags.
type Bookmark struct {
	Model

	Name  string
	URL   string
	Color string

	OwnerID uint
	Folders []Folder `gorm:"many2many:bookmark_folder;"`
	Tags    []Tag    `gorm:"many2many:bookmark_tag;"`
}

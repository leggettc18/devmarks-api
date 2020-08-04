package model

type Bookmark struct {
	Model

	Name  string
	Url   string
	Color string

	OwnerID uint
	Folders []Folder `gorm:"many2many:bookmark_folder;"`
	Tags    []Tag    `gorm:"many2many:bookmark_tag;"`
}

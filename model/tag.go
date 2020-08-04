package model

type Tag struct {
	Model

	Name string

	Bookmarks []Bookmark `gorm:"many2many:bookmark_tag;"`
}

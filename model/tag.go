package model

// Tag is a model representing tags that our app can save. They can be used to group any number of
// bookmarks together.
type Tag struct {
	Model

	Name string

	Bookmarks []Bookmark `gorm:"many2many:bookmark_tag;"`
}

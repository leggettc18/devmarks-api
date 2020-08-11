package db

import (
	"github.com/pkg/errors"

	"leggett.dev/devmarks/api/model"
)

// GetBookmarkByID queries the database for a bookmark with the specified id
func (db *Database) GetBookmarkByID(id uint) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	return &bookmark, errors.Wrap(db.First(&bookmark, id).Error, "unable to get bookmark")
}

// GetBookmarksByUserID returns all the bookmarks from the database that are owned by the user
// corresponding to the userID provided.
func (db *Database) GetBookmarksByUserID(userID uint) ([]*model.Bookmark, error) {
	var bookmarks []*model.Bookmark
	return bookmarks, errors.Wrap(db.Find(&bookmarks, model.Bookmark{OwnerID: userID}).Error, "unable to get bookmarks")
}

// CreateBookmark inserts the specified bookmark into the database.
func (db *Database) CreateBookmark(bookmark *model.Bookmark) error {
	return errors.Wrap(db.Create(bookmark).Error, "unable to create bookmark")
}

// UpdateBookmark updates the specified bookmark in the database.
func (db *Database) UpdateBookmark(bookmark *model.Bookmark) error {
	return errors.Wrap(db.Save(bookmark).Error, "unable to update bookmark")
}

// DeleteBookmarkByID deletes the bookmark with the specified ID from the databse.
func (db *Database) DeleteBookmarkByID(id uint) error {
	return errors.Wrap(db.Delete(&model.Bookmark{}, id).Error, "unable to delete todo")
}

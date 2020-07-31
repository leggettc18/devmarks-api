package db

import (
	"github.com/pkg/errors"

	"leggett.dev/devmarks/api/model"
)

func (db *Database) GetBookmarkById(id uint) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	return &bookmark, errors.Wrap(db.First(&bookmark, id).Error, "unable to get bookmark")
}

func (db *Database) GetBookmarksByUserId(userId uint) ([]*model.Bookmark, error) {
	var bookmarks []*model.Bookmark
	return bookmarks, errors.Wrap(db.Find(&bookmarks, model.Bookmark{UserID: userId}).Error, "unable to get bookmarks")
}

func (db *Database) CreateBookmark(bookmark *model.Bookmark) error {
	return errors.Wrap(db.Create(bookmark).Error, "unable to create bookmark")
}

func (db *Database) UpdateBookmark(bookmark *model.Bookmark) error {
	return errors.Wrap(db.Save(bookmark).Error, "unable to update bookmark")
}

func (db *Database) DeleteBookmarkById(id uint) error {
	return errors.Wrap(db.Delete(&model.Bookmark{}, id).Error, "unable to delete todo")
}

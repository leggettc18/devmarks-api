package app

import "leggett.dev/devmarks/api/model"

// GetBookmarkByID returns a Bookmark model from the bookmark's ID
func (ctx *Context) GetBookmarkByID(id uint) (*model.Bookmark, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}

	bookmark, err := ctx.Database.GetBookmarkByID(id)
	if err != nil {
		return nil, err
	}

	if bookmark.OwnerID != ctx.User.ID {
		return nil, ctx.AuthorizationError()
	}

	return bookmark, nil
}

func (ctx *Context) getBookmarksByUserID(userID uint) ([]*model.Bookmark, error) {
	return ctx.Database.GetBookmarksByUserID(userID)
}

// GetUserBookmarks returns a slice of all bookmark models that are owned by
// The currently authenticated User
func (ctx *Context) GetUserBookmarks() ([]*model.Bookmark, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}

	return ctx.getBookmarksByUserID(ctx.User.ID)
}

// CreateBookmark performs the business logic necessary to create and
// validate a Bookmark given an initial instance of one
func (ctx *Context) CreateBookmark(bookmark *model.Bookmark) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	bookmark.OwnerID = ctx.User.ID

	if err := ctx.validateBookmark(bookmark); err != nil {
		return err
	}

	return ctx.Database.CreateBookmark(bookmark)
}

const maxBookmarkNameLength = 100

func (ctx *Context) validateBookmark(bookmark *model.Bookmark) *ValidationError {
	if len(bookmark.Name) > maxBookmarkNameLength {
		return &ValidationError{"name is too long"}
	}

	return nil
}

// UpdateBookmark performs the business logic necessary to validate and update
// a given bookmark model
func (ctx *Context) UpdateBookmark(bookmark *model.Bookmark) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	if bookmark.OwnerID != ctx.User.ID {
		return ctx.AuthorizationError()
	}

	if bookmark.ID == 0 {
		return &ValidationError{"cannot update"}
	}

	if err := ctx.validateBookmark(bookmark); err != nil {
		return nil
	}

	return ctx.Database.UpdateBookmark(bookmark)
}

// DeleteBookmarkByID performs the necessary business logic to delete
// a bookmark owned by the currently authenticated user
func (ctx *Context) DeleteBookmarkByID(id uint) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	bookmark, err := ctx.GetBookmarkByID(id)
	if err != nil {
		return err
	}

	if bookmark.OwnerID != ctx.User.ID {
		return ctx.AuthorizationError()
	}

	return ctx.Database.DeleteBookmarkByID(id)
}

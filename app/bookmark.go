package app

import "leggett.dev/devmarks/api/model"

func (ctx *Context) GetBookmarkById(id uint) (*model.Bookmark, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}

	bookmark, err := ctx.Database.GetBookmarkById(id)
	if err != nil {
		return nil, err
	}

	if bookmark.OwnerID != ctx.User.ID {
		return nil, ctx.AuthorizationError()
	}

	return bookmark, nil
}

func (ctx *Context) getBookmarksByUserId(userId uint) ([]*model.Bookmark, error) {
	return ctx.Database.GetBookmarksByUserId(userId)
}

func (ctx *Context) GetUserBookmarks() ([]*model.Bookmark, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}

	return ctx.getBookmarksByUserId(ctx.User.ID)
}

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

func (ctx *Context) DeleteBookmarkById(id uint) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	bookmark, err := ctx.GetBookmarkById(id)
	if err != nil {
		return err
	}

	if bookmark.OwnerID != ctx.User.ID {
		return ctx.AuthorizationError()
	}

	return ctx.Database.DeleteBookmarkById(id)
}

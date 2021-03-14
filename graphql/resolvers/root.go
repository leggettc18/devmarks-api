package resolvers

import (
	"context"
	"errors"

	"leggett.dev/devmarks/api/db"
	"leggett.dev/devmarks/api/model"
)

type RootResolver struct {
	Database *db.Database
}

func NewRoot(db *db.Database) (*RootResolver, error) {
	r := &RootResolver{Database: db}
	return r, nil
}

func(r RootResolver) Bookmarks(ctx context.Context) (*[]*BookmarkResolver, error){
	var user, ok = ctx.Value("user").(*model.User)

	if !ok {
		return nil, errors.New("bookmarks: no authenticated user in context")
	}
	bookmarks, err := r.Database.GetBookmarksByUserID(user.ID)

	if err != nil {
		return nil, err
	}

	var resolvers []*BookmarkResolver
	for _, bookmark := range bookmarks {
		resolvers = append(resolvers, &BookmarkResolver{*bookmark})
	}
	return &resolvers, nil
}

type NewBookmarkArgs struct {
	Name string
	Url string
	Color *string
}

func (r RootResolver) NewBookmark(ctx context.Context, args NewBookmarkArgs) (*BookmarkResolver, error) {
	var user, ok = ctx.Value("user").(*model.User)

	if !ok {
		return nil, errors.New("bookmarks: no authenticated user in context")
	}

	newBookmark := model.Bookmark {
		Name: args.Name,
		URL: args.Url,
		Color: *args.Color,
		OwnerID: user.ID,
	}

	if err := r.Database.CreateBookmark(&newBookmark); err != nil {
		return nil, err
	}

	bookmarkResolver := &BookmarkResolver{newBookmark}

	return bookmarkResolver, nil
}
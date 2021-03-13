package resolvers

import (
	"context"

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
	var user *model.User = ctx.Value("user").(*model.User)
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
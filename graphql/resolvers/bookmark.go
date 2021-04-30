package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"leggett.dev/devmarks/api/db"
	"leggett.dev/devmarks/api/graphql/helpers"
	"leggett.dev/devmarks/api/model"
)

type BookmarkResolver struct {
	Bookmark model.Bookmark
	DB db.Database
}

func (r *BookmarkResolver) ID() graphql.ID {
	return *helpers.GqlIDP(r.Bookmark.ID)
}

func (r *BookmarkResolver) CreatedAt() graphql.Time {
	return graphql.Time{ Time: r.Bookmark.CreatedAt }
}

func (r *BookmarkResolver) UpdatedAt() graphql.Time {
	return graphql.Time{ Time: r.Bookmark.UpdatedAt }
}

func (r *BookmarkResolver) Name() string {
	return r.Bookmark.Name
}

func (r *BookmarkResolver) URL() string {
	return r.Bookmark.URL
}

func (r *BookmarkResolver) Color() *string {
	return r.Bookmark.Color
}

func (r *BookmarkResolver) Owner() (*UserResolver, error) {
	user, err := r.DB.GetUserById(r.Bookmark.OwnerID)
	if err != nil {
		return nil, err
	}
	return &UserResolver{ *user }, nil
}

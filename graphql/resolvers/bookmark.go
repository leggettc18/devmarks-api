package resolvers

import (
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"leggett.dev/devmarks/api/model"
)

type BookmarkResolver struct {
	Bookmark model.Bookmark
}

func (r *BookmarkResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.Bookmark.ID))
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

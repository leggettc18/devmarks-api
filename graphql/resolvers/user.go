package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"leggett.dev/devmarks/api/graphql/helpers"
	"leggett.dev/devmarks/api/model"
)

type UserResolver struct {
	User model.User
}

func (r *UserResolver) ID() graphql.ID {
	return *helpers.GqlIDP(r.User.ID)
}

func (r *UserResolver) CreatedAt() graphql.Time {
	return graphql.Time{ Time: r.User.CreatedAt }
}

func (r *UserResolver) UpdatedAt() graphql.Time {
	return graphql.Time{ Time: r.User.UpdatedAt }
}

func (r *UserResolver) Email() string {
	return r.User.Email
}

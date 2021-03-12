package resolvers

import "leggett.dev/devmarks/api/app"

type RootResolver struct {
	ctx app.Context
}

func NewRoot(ctx app.Context) (*RootResolver, error) {
	r := &RootResolver{}
	return r, nil
}

func(r *RootResolver) Bookmarks() (*[]*BookmarkResolver, error){
	bookmarks, err := r.ctx.Database.GetBookmarksByUserID(r.ctx.User.ID)

	if err != nil {
		return nil, err
	}

	var resolvers []*BookmarkResolver
	for _, bookmark := range bookmarks {
		resolvers = append(resolvers, &BookmarkResolver{*bookmark})
	}
	return &resolvers, nil
}
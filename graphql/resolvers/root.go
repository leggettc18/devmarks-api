package resolvers

import (
	"context"
	"errors"

	"leggett.dev/devmarks/api/app"
	"leggett.dev/devmarks/api/model"
)

type RootResolver struct {
	App *app.App
}

func NewRoot(app *app.App) (*RootResolver, error) {
	r := &RootResolver{App: app}
	return r, nil
}

func(r RootResolver) Bookmarks(ctx context.Context) (*[]*BookmarkResolver, error){
	var user, ok = ctx.Value("user").(*model.User)

	if !ok {
		return nil, errors.New("bookmarks: no authenticated user in context")
	}
	bookmarks, err := r.App.Database.GetBookmarksByUserID(user.ID)

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

	if err := r.App.Database.CreateBookmark(&newBookmark); err != nil {
		return nil, err
	}

	bookmarkResolver := &BookmarkResolver{newBookmark}

	return bookmarkResolver, nil
}

type LoginArgs struct {
	Email	 string
	Password string
}

func (r *RootResolver) Login(args LoginArgs) (*AuthResolver, error) {
	user, errUser := r.App.Database.GetUserByEmail(args.Email)
	if errUser != nil {
		return nil, errUser
	}
	model.ComparePasswordHash(user.HashedPassword, []byte(args.Password))

	token, errToken := GenerateToken(user, r.App.Config.SecretKey)
	if errToken != nil {
		return nil, errToken
	}
	payload := AuthPayload{
		Token: &token,
		User: user,
	}
	return &AuthResolver{payload}, nil
}
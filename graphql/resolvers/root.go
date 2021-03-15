package resolvers

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
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

type RegisterArgs struct {
	Email string
	Password string
}

func (r *RootResolver) Register(args RegisterArgs) (*AuthResolver, error) {
	passwordHash, errHash := bcrypt.GenerateFromPassword(
		[]byte(args.Password),
		bcrypt.DefaultCost,
	)
	if errHash != nil {
		return nil, errHash
	}

	newUser := model.User{
		Email: args.Email,
		HashedPassword: passwordHash,
	}

	if err := r.App.Database.CreateUser(&newUser); err != nil {
		return nil, err
	}

	token, errToken := GenerateToken(&newUser, r.App.Config.SecretKey)
	if errToken != nil {
		return nil, errToken
	}

	payload := AuthPayload{
		Token: &token,
		User: &newUser,
	}

	return &AuthResolver{payload}, nil
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
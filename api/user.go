package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/bearer"
	"leggett.dev/devmarks/api/app"
	"leggett.dev/devmarks/api/model"
)

// UserInput represents the input to the CreateUser function
type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse represents the response written to the HTTP response header upon
// CreateUser's completion
type UserResponse struct {
	ID uint `json:"id"`
}

// CreateUser creates a new user based on the json data provided in the HTTP Request
func (a *API) CreateUser(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input UserInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	user := &model.User{Email: input.Email}

	if err := ctx.CreateUser(user, input.Password); err != nil {
		return err
	}

	data, err := json.Marshal(&UserResponse{ID: user.ID})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) validateLogin(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	user, err := a.App.GetUserByEmail(userName)

	if user == nil || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	if ok := user.CheckPassword(password); !ok {
		return nil, errors.Wrap(err, "invalid credentials")
	}
	return auth.NewDefaultUser(user.Email, strconv.Itoa(int(user.ID)), nil, nil), nil
}

func (a *API) createToken(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}
	token := uuid.New().String()
	user := auth.NewDefaultUser(ctx.User.Email, strconv.Itoa(int(ctx.User.ID)), nil, nil)
	tokenStrategy := authenticator.Strategy(bearer.CachedStrategyKey)
	auth.Append(tokenStrategy, token, user, r)
	body := fmt.Sprintf("token: %s \n", token)
	_, err := w.Write([]byte(body))
	return err
}

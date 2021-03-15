package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shaj13/go-guardian/auth"
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

// GetUser Retrieves the authenticated user from the database
func (a *API) GetUser(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	user := ctx.User

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) validateLogin(r *http.Request, userName, password string) (auth.Info, error) {
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
	var input UserInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	user, err := a.validateLogin(r, input.Email, input.Password)
	if err != nil {
		return err
	}

	bearerToken := uuid.New().String()
	a.App.AuthCache.Store(bearerToken, user, r)
	responseBody := fmt.Sprintf("{ \"token\": \"%s\" }\n", bearerToken)
	_, err = w.Write([]byte(responseBody))
	return err
}

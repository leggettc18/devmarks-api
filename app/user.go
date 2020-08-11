package app

import (
	"strings"

	"github.com/pkg/errors"

	"leggett.dev/devmarks/api/model"
)

// GetUserByEmail returns a user that matches the specified email address
func (a *App) GetUserByEmail(email string) (*model.User, error) {
	return a.Database.GetUserByEmail(email)
}

// CreateUser performs the business logic necessary to create and validate a new
// User, returning an error if validation fails or a password cannot be set.
func (ctx *Context) CreateUser(user *model.User, password string) error {
	if err := ctx.validateUser(user, password); err != nil {
		return err
	}

	if err := user.SetPassword(password); err != nil {
		return errors.Wrap(err, "unable to set user password")
	}

	return ctx.Database.CreateUser(user)
}

func (ctx *Context) validateUser(user *model.User, password string) *ValidationError {
	// naive email validation
	if !strings.Contains(user.Email, "@") {
		return &ValidationError{"invalid email"}
	}

	if password == "" {
		return &ValidationError{"password is required"}
	}

	return nil
}

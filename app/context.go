package app

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"leggett.dev/devmarks/api/db"
	"leggett.dev/devmarks/api/model"
)

// Context represents the current Context of our application (logger, remote address,
// currently logged in user, etc.))
type Context struct {
	Logger        logrus.FieldLogger
	RemoteAddress string
	Database      *db.Database
	User          *model.User
}

// WithLogger returns an instance of the context it was called on with the specified logger
// substituted in.
func (ctx *Context) WithLogger(logger logrus.FieldLogger) *Context {
	ret := *ctx
	ret.Logger = logger
	return &ret
}

// WithRemoteAddress returns an instance of the context it was called on with the specified
// Remote Address substituted in.
func (ctx *Context) WithRemoteAddress(address string) *Context {
	ret := *ctx
	ret.RemoteAddress = address
	return &ret
}

// WithUser returns an instance of the context it was called on with the specified
// User model substituted in.
func (ctx *Context) WithUser(user *model.User) *Context {
	ret := *ctx
	ret.User = user
	return &ret
}

// AuthorizationError returns a UserError signifying a failed authorization
func (ctx *Context) AuthorizationError() *UserError {
	return &UserError{Message: "unauthorized", StatusCode: http.StatusForbidden}
}

package app

import (
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/store"
	"github.com/sirupsen/logrus"
	"leggett.dev/devmarks/api/db"
)

// App is an object representing our App's configuration
type App struct {
	Config   *Config
	Database *db.Database
	AuthCache store.Cache
	Authenticator auth.Authenticator
}

// NewContext returns a new Context object
func (a *App) NewContext() *Context {
	return &Context{
		Logger:   logrus.New(),
		Database: a.Database,
	}
}

// New returns a new App object
func New() (app *App, err error) {
	app = &App{}
	app.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}
	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}
	app.Database, err = db.New(dbConfig)
	if err != nil {
		return nil, err
	}
	return app, err
}

// Close performs any actions necessary to close our our running
// app, like closing the database connection
func (a *App) Close() error {
	return a.Database.Close()
}

// ValidationError contains specific information about why a validation
// failure occurred.
type ValidationError struct {
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// UserError contains specific information about what User-related
// error occurred and what status code to write to the HTTP response
// header
type UserError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *UserError) Error() string {
	return e.Message
}

package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/friendsofgo/graphiql"
	"github.com/gorilla/mux"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/token"
	"github.com/shaj13/go-guardian/store"
	"github.com/sirupsen/logrus"

	"leggett.dev/devmarks/api/app"
	"leggett.dev/devmarks/api/graphql/resolvers"
	"leggett.dev/devmarks/api/model"
)

var (
	opts = []graphql.SchemaOpt{graphql.UseStringDescriptions()}
)


type statusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// API is an object representing our API's configuration, and includes a pointer
// to our App's App object
type API struct {
	App    *app.App
	Config *Config
}

// New returns a new API object from our App's App object
func New(a *app.App) (api *API, err error) {
	api = &API{App: a}
	api.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}
	return api, nil
}

func (a *API) setupGoGuardian() {
	a.App.Authenticator = auth.New()
	a.App.AuthCache = store.NewFIFO(context.Background(), time.Minute*10)

	tokenStrategy := token.New(token.NoOpAuthenticate, a.App.AuthCache)

	a.App.Authenticator.EnableStrategy(token.CachedStrategyKey, tokenStrategy)
}

// Init Initializes our API (routes, authentication setup, etc.)
func (a *API) Init(r *mux.Router) {
	// authentication
	a.setupGoGuardian()
	r.Handle("/auth/token/", a.handler(a.createToken)).Methods("POST")

	// user methods
	r.Handle("/users/", a.handler(a.CreateUser)).Methods("POST")
	r.Handle("/user/", a.handler(a.GetUser)).Methods("GET")

	// bookmark methods
	bookmarksRouter := r.PathPrefix("/bookmarks").Subrouter()
	bookmarksRouter.Handle("/", a.handler(a.GetBookmarks)).Methods("GET")
	bookmarksRouter.Handle("/", a.handler(a.CreateBookmark)).Methods("POST")
	bookmarksRouter.Handle("/{id:[0-9]+}/", a.handler(a.GetBookmarkByID)).Methods("GET")
	bookmarksRouter.Handle("/{id:[0-9]+}/", a.handler(a.UpdateBookmarkByID)).Methods("PATCH")
	bookmarksRouter.Handle("/{id:[0-9]+}/", a.handler(a.DeleteBookmarkByID)).Methods("DELETE")
}

func parseSchema(path string, resolver interface{}) *graphql.Schema {
	bstr, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	schemaString := string(bstr)
	parsedSchema, err := graphql.ParseSchema(
		schemaString,
		resolver,
		opts...,
	)
	if err != nil {
		panic(err)
	}
	return parsedSchema
}

func (a *API) InitGraphql(r *mux.Router) {
	// graphql
	rootResolver, err := resolvers.NewRoot(a.App)
	schema := parseSchema("./schema.graphql", rootResolver)
	wsHandler:= graphqlws.NewHandlerFunc(
		schema,
		&relay.Handler{
			Schema: schema,
		},
	)
	r.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)
		beginTime := time.Now()
		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker: hijacker,
		}
		token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
		ctx := context.WithValue(context.Background(), "token", token)
		ctx = context.WithValue(ctx, "remote_address", a.IPAddressForRequest((r)))
		logger := logrus.New()
		ctx = context.WithValue(ctx, "logger", logger.WithField("request_id", base64.RawURLEncoding.EncodeToString(model.NewID())))
		wsHandler.ServeHTTP(w, r.WithContext(ctx))
		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)

			logger := logger.WithFields(logrus.Fields{
				"duration":    duration,
				"status_code": statusCode,
				"remote":      ctx.Value("remote_address").(string),
			})
			logger.Info(r.Method + " " + r.URL.RequestURI())
		}()

		defer func() {
			if r := recover(); r != nil {
				logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		w.Header().Set("Content-Type", "application/json")
	})
	// graphiql
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	if err != nil {
		panic(err)
	}
	r.Handle("/graphiql", graphiqlHandler).Methods("GET")
}

func (a *API) handler(f func(*app.Context, http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)

		beginTime := time.Now()

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		ctx := a.App.NewContext().WithRemoteAddress(a.IPAddressForRequest(r))
		ctx = ctx.WithLogger(ctx.Logger.WithField("request_id", base64.RawURLEncoding.EncodeToString(model.NewID())))

		/* if username, password, ok := r.BasicAuth(); ok {
			user, err := a.App.GetUserByEmail(username)

			if user == nil || err != nil {
				if err != nil {
					ctx.Logger.WithError(err).Error("unable to get user")
				}
				http.Error(w, "invalid credentials", http.StatusForbidden)
				return
			}

			if ok := user.CheckPassword(password); !ok {
				http.Error(w, "invalid credentials", http.StatusForbidden)
			}

			ctx = ctx.WithUser(user)
		} */
		if !(r.URL.Path == "/api/users/" || r.URL.Path == "/api/auth/token/") {
			tokenStrategy := a.App.Authenticator.Strategy(token.CachedStrategyKey)
			userInfo, err := tokenStrategy.Authenticate(r.Context(), r)
			if err != nil {
				ctx.Logger.WithError(err).Error("unable to get user")
				http.Error(w, "invalid credentials", http.StatusForbidden)
				return
			}
			user, err := a.App.GetUserByEmail(userInfo.UserName())

			if user == nil || err != nil {
				if err != nil {
					ctx.Logger.WithError(err).Error("unable to get user")
				}
				http.Error(w, "invalid credentials", http.StatusForbidden)
				return
			}

			ctx = ctx.WithUser(user)
		}

		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)

			logger := ctx.Logger.WithFields(logrus.Fields{
				"duration":    duration,
				"status_code": statusCode,
				"remote":      ctx.RemoteAddress,
			})
			logger.Info(r.Method + " " + r.URL.RequestURI())
		}()

		defer func() {
			if r := recover(); r != nil {
				ctx.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		w.Header().Set("Content-Type", "application/json")

		if err := f(ctx, w, r); err != nil {
			if verr, ok := err.(*app.ValidationError); ok {
				data, err := json.Marshal(verr)
				if err == nil {
					w.WriteHeader(http.StatusBadRequest)
					_, err = w.Write(data)
				}

				if err != nil {
					ctx.Logger.Error(err)
					http.Error(w, "interval server error", http.StatusInternalServerError)
				}
			} else if uerr, ok := err.(*app.UserError); ok {
				data, err := json.Marshal(uerr)
				if err == nil {
					w.WriteHeader(uerr.StatusCode)
					_, err = w.Write(data)
				}

				if err != nil {
					ctx.Logger.Error(err)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			} else {
				ctx.Logger.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}
	})
}

// IPAddressForRequest gets the IP address from our HTTP request
func (a *API) IPAddressForRequest(r *http.Request) string {
	addr := r.RemoteAddr
	if a.Config.ProxyCount > 0 {
		h := r.Header.Get("X-Forwarded-For")
		if h != "" {
			clients := strings.Split(h, ",")
			if a.Config.ProxyCount > len(clients) {
				addr = clients[0]
			} else {
				addr = clients[len(clients)-a.Config.ProxyCount]
			}
		}
	}
	//TODO: consider refactoring to use regex instead.
	if (strings.Contains(addr, "[")) { //If addr is ipv6
		sep_strings := strings.Split(strings.TrimSpace(addr), ":") //split string at the colons
		sep_strings = sep_strings[:len(sep_strings) -1] //remove the last string (the port number)
		return strings.Join(sep_strings, ":") //Join the remaining elements back together into one string with the colons in between.
	} //If addr is ipv4
	return strings.Split(strings.TrimSpace(addr), ":")[0]
}

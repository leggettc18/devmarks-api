package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"leggett.dev/devmarks/api/app"
	"leggett.dev/devmarks/api/model"
)

// GetBookmarks returns the bookmarks corresponding to the currently authenticated user in json form
func (a *API) GetBookmarks(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	bookmarks, err := ctx.GetUserBookmarks()
	if err != nil {
		return err
	}

	data, err := json.Marshal(bookmarks)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// CreateBookmarkInput represents the input to the CreateBookmark function
type CreateBookmarkInput struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Color string `json:"color"`
}

// CreateBookmarkResponse represents the response that will be sent upon completion of
// the CreateBookmark function
type CreateBookmarkResponse struct {
	ID uint `json:"id"`
}

// CreateBookmark creates a new bookmark owned by the currently authenticated user based
// on json from the HTTP Request
func (a *API) CreateBookmark(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input CreateBookmarkInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	bookmark := &model.Bookmark{Name: input.Name, URL: input.URL, Color: &input.Color}

	if err := ctx.CreateBookmark(bookmark); err != nil {
		return err
	}

	data, err := json.Marshal(&CreateBookmarkResponse{ID: bookmark.ID})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// GetBookmarkByID writes the json representation of a bookmark to the HTTP Response Header,
// if the currently authenticated user has access to it.
func (a *API) GetBookmarkByID(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIDFromRequest(r)
	bookmark, err := ctx.GetBookmarkByID(id)
	if err != nil {
		return err
	}

	data, err := json.Marshal(bookmark)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// UpdateBookmarkInput represents the input to the UpdateBookmark function
type UpdateBookmarkInput struct {
	Name  *string `json:"name"`
	URL   *string `json:"url"`
	Color *string `json:"color"`
}

// UpdateBookmarkByID updates the bookmark whose ID is specified in the HTTP request if it is owned
// by the currently authenticated user.
func (a *API) UpdateBookmarkByID(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIDFromRequest(r)

	var input UpdateBookmarkInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	existingBookmark, err := ctx.GetBookmarkByID(id)
	if err != nil || existingBookmark == nil {
		return err
	}

	if input.Name != nil {
		existingBookmark.Name = *input.Name
	}
	if input.URL != nil {
		existingBookmark.URL = *input.URL
	}
	if input.Color != nil {
		existingBookmark.Color = input.Color
	}

	err = ctx.UpdateBookmark(existingBookmark)
	if err != nil {
		return err
	}

	data, err := json.Marshal(existingBookmark)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// DeleteBookmarkByID deletes the bookmark whose ID is specified in the HTTP request if it is
// owned by the currently authenticated user
func (a *API) DeleteBookmarkByID(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIDFromRequest(r)
	err := ctx.DeleteBookmarkByID(id)

	if err != nil {
		return err
	}

	return &app.UserError{StatusCode: http.StatusOK, Message: "removed"}
}

func getIDFromRequest(r *http.Request) uint {
	vars := mux.Vars(r)
	id := vars["id"]

	intID, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return 0
	}

	return uint(intID)
}

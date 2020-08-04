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

type CreateBookmarkInput struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Color string `json:"color"`
}

type CreateBookmarkResponse struct {
	Id uint `json:"id"`
}

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

	bookmark := &model.Bookmark{Name: input.Name, Url: input.Url, Color: input.Color}

	if err := ctx.CreateBookmark(bookmark); err != nil {
		return err
	}

	data, err := json.Marshal(&CreateBookmarkResponse{Id: bookmark.ID})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) GetBookmarkById(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIdFromRequest(r)
	bookmark, err := ctx.GetBookmarkById(id)
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

type UpdateBookmarkInput struct {
	Name  *string `json:"name"`
	Url   *string `json:"url`
	Color *string `json:"color"`
}

func (a *API) UpdateBookmarkById(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIdFromRequest(r)

	var input UpdateBookmarkInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	existingBookmark, err := ctx.GetBookmarkById(id)
	if err != nil || existingBookmark == nil {
		return err
	}

	if input.Name != nil {
		existingBookmark.Name = *input.Name
	}
	if input.Url != nil {
		existingBookmark.Url = *input.Url
	}
	if input.Color != nil {
		existingBookmark.Color = *input.Color
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

func (a *API) DeleteBookmarkById(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIdFromRequest(r)
	err := ctx.DeleteBookmarkById(id)

	if err != nil {
		return err
	}

	return &app.UserError{StatusCode: http.StatusOK, Message: "removed"}
}

func getIdFromRequest(r *http.Request) uint {
	vars := mux.Vars(r)
	id := vars["id"]

	intId, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return 0
	}

	return uint(intId)
}

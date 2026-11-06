package handler

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest/view"
)

// UserHandler handles HTTP requests for user resources.
type UserHandler struct {
	App         *gorest.App
	UserLibrary *library.UserLibrary
	JSON        *view.JSON
}

// NewUserHandler creates a handler connected to the application and user library.
func NewUserHandler(app *gorest.App, lib *library.UserLibrary) *UserHandler {
	return &UserHandler{
		App:         app,
		UserLibrary: lib,
	}
}

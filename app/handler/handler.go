package handler

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/view"
)

// UserHandler handles HTTP requests for user resources.
type UserHandler struct {
	App         *gorest.App
	UserLibrary UserUseCase
	JSON        *view.JSON
}

type UserUseCase interface {
	Login(*gorest.Context, *entity.User) (string, int, error)
	CreateUser(*entity.User) (string, int, error)
	GetUsers() (*[]entity.User, string, int, error)
}

// NewUserHandler creates a handler connected to the application and user library.
func NewUserHandler(app *gorest.App, lib UserUseCase) *UserHandler {
	return &UserHandler{
		App:         app,
		UserLibrary: lib,
	}
}

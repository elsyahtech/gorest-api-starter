package library

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/model"
)

// UserLibrary contains user-related application logic.
type UserLibrary struct {
	App       *gorest.App
	UserModel *model.UserModel
}

// NewUserLibrary creates a user library backed by the given model.
func NewUserLibrary(app *gorest.App, userModel *model.UserModel) *UserLibrary {
	return &UserLibrary{
		App:       app,
		UserModel: userModel,
	}
}

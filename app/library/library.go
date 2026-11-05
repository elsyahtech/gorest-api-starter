package library

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// UserLibrary contains user-related application logic.
type UserLibrary struct {
	App       *gorest.App
	UserModel UserRepository
}

// NewUserLibrary creates a user library backed by the given model.
func NewUserLibrary(app *gorest.App, userModel UserRepository) *UserLibrary {
	return &UserLibrary{
		App:       app,
		UserModel: userModel,
	}
}

type UserRepository interface {
	CreateUser(*entity.User) (string, int, error)
	GetUserByEmail(string) (*entity.User, string, int, error)
	GetUsers() (*[]entity.User, string, int, error)
}

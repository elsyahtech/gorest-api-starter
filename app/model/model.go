package model

import "github.com/elsyahtech/gorest"

// UserModel provides database operations for users.
type UserModel struct {
	App *gorest.App
}

// NewUserModel creates a user model backed by the given application.
func NewUserModel(app *gorest.App) *UserModel {
	return &UserModel{
		App: app,
	}
}

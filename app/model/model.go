package model

import "github.com/elsyahtech/gorest"

// UserModel provides database operations for users.
type UserModel struct {
	App *gorest.App
}

// DepartmentModel provides database operations for departments.
type DepartmentModel struct {
	App *gorest.App
}

// OrderModel provides database operations for orders.
type OrderModel struct {
	App *gorest.App
}

// NewUserModel creates a user model backed by the given application.
func NewUserModel(app *gorest.App) *UserModel {
	return &UserModel{
		App: app,
	}
}

// NewDepartmentModel creates a department model backed by the given application.
func NewDepartmentModel(app *gorest.App) *DepartmentModel {
	return &DepartmentModel{
		App: app,
	}
}

// NewOrderModel creates an order model backed by the given application.
func NewOrderModel(app *gorest.App) *OrderModel {
	return &OrderModel{
		App: app,
	}
}

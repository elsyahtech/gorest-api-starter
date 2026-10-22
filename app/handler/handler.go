package handler

import (
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest/view"
)

// UserHandler handles HTTP requests for user resources.
type UserHandler struct {
	App         *gorest.App
	UserLibrary *library.UserLibrary
	JSON        *view.JSON
}

// DepartmentHandler handles HTTP requests for department resources.
type DepartmentHandler struct {
	App               *gorest.App
	DepartmentLibrary *library.DepartmentLibrary
	JSON              *view.JSON
}

// OrderHandler handles HTTP requests for order resources.
type OrderHandler struct {
	App          *gorest.App
	OrderLibrary *library.OrderLibrary
	JSON         *view.JSON
}

// NewUserHandler creates a handler connected to the application and user library.
func NewUserHandler(app *gorest.App, lib *library.UserLibrary) *UserHandler {
	return &UserHandler{
		App:         app,
		UserLibrary: lib,
	}
}

// NewDepartmentHandler creates a handler connected to the application and department library.
func NewDepartmentHandler(app *gorest.App, departmentLibrary *library.DepartmentLibrary) *DepartmentHandler {
	return &DepartmentHandler{
		App:               app,
		DepartmentLibrary: departmentLibrary,
	}
}

// NewOrderHandler creates a handler connected to the application and order library.
func NewOrderHandler(app *gorest.App, orderLibrary *library.OrderLibrary) *OrderHandler {
	return &OrderHandler{
		App:          app,
		OrderLibrary: orderLibrary,
	}
}

package router

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest/server"
)

// UserRouter registers runnable user-query examples under the /examples prefix.
// Static paths are registered before parameter paths such as /users/:id.
func UserRouter(app *gorest.App) {
	mdl := model.NewUserModel(app)
	lib := library.NewUserLibrary(app, mdl)
	hndlr := handler.NewUserHandler(app, lib)
	RegisterUserRoutes(app, app.Group("/users"), hndlr)
}

// RegisterUserRoutes registers user endpoints on a route group.
func RegisterUserRoutes(app *gorest.App, group *server.Server, hndlr *handler.UserHandler) {
	authentication := app.Filter().Authentication()

	group.Get("/", authentication, hndlr.GetUsers)
	group.Post("/create", hndlr.CreateUser)
	group.Post("/login", hndlr.Login)
}

package router

import (
	"net/http"

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
	lib := library.NewUserLibrary(mdl)
	hndlr := handler.NewUserHandler(app, lib)
	RegisterUserRoutes(app.Group("/examples"), hndlr)
}

// RegisterUserRoutes registers user endpoints on a route group.
func RegisterUserRoutes(group *server.Server, hndlr *handler.UserHandler) {
	// Read users with ORM queries.
	group.Get("/orm-queries/users", hndlr.GetUsers)
	group.Get("/orm-queries/users/select-all", hndlr.GetUsersWithAllColumns)
	group.Get("/orm-queries/users/select-columns", hndlr.GetUsersWithSelectedColumns)
	group.Get("/orm-queries/users/order-by", hndlr.GetUsersOrderedByID)
	group.Get("/orm-queries/users/group-by", hndlr.GetUserCountsByDepartment)
	group.Get("/orm-queries/users/distinct", hndlr.GetUsersWithDistinctDepartmentIDs)
	group.Get("/orm-queries/users/having", hndlr.GetDepartmentsByUserCount)
	group.Get("/orm-queries/users/join", hndlr.GetUsersWithDepartment)
	group.Get("/orm-queries/users/preload", hndlr.GetUsersWithPreloadedDepartment)
	group.Get("/orm-queries/users/limit", hndlr.GetUsersByLimit)
	group.Get("/orm-queries/users/pagination", hndlr.GetUsersByLimitOffset)

	// Read users with each native query API. The SQL route uses the configured SQL driver.
	group.Get("/native-queries/users/sql", hndlr.GetUsersWithNativeSQL)
	group.Get("/native-queries/users/mongodb", hndlr.GetUsersWithNativeMongoQuery)
	group.Get("/native-queries/users/scylla", hndlr.GetUsersWithNativeCQL)

	// HEAD exposes the GET resource metadata without returning a response body.
	group.Head("/orm-queries/users", hndlr.GetUsers)
	group.Head("/orm-queries/users/:id", hndlr.GetUserByID)

	// OPTIONS advertises methods available on the collection and a single user.
	group.Options("/orm-queries/users", func(ctx *gorest.Context) error {
		ctx.Set("Allow", "GET, HEAD, POST, OPTIONS")
		return ctx.SendStatus(http.StatusNoContent)
	})
	group.Options("/orm-queries/users/:id", func(ctx *gorest.Context) error {
		ctx.Set("Allow", "GET, HEAD, PUT, PATCH, DELETE, OPTIONS")
		return ctx.SendStatus(http.StatusNoContent)
	})

	// Keep this parameter route after the static GET routes above.
	group.Get("/orm-queries/users/:id", hndlr.GetUserByID)

	// Create users with ORM features.
	// Some transaction and bulk-write examples call BeginTx. GOREST's
	// BeginTx/Commit/Rollback calls are no-ops for MongoDB and ScyllaDB.
	group.Post("/orm-queries/users", hndlr.CreateUser)
	group.Post("/orm-queries/users/transaction", hndlr.CreateUserWithTransaction)
	group.Post("/orm-queries/users/selected-columns", hndlr.CreateUserWithSelectedColumns)
	group.Post("/orm-queries/users/upsert", hndlr.UpsertUser)
	group.Post("/orm-queries/users/returning", hndlr.CreateUserReturningResult)
	group.Post("/orm-queries/users/bulk", hndlr.CreateUsersInBulk)

	// Create users with native query APIs.
	group.Post("/native-queries/users/sql", hndlr.CreateUserWithNativeSQL)
	group.Post("/native-queries/users/mongodb", hndlr.CreateUserWithNativeMongoQuery)
	group.Post("/native-queries/users/scylla", hndlr.CreateUserWithNativeCQL)

	// Update users. Register fixed paths before the /:id parameter route.
	group.Put("/orm-queries/users/bulk", hndlr.UpdateUsersInBulk)
	group.Put("/orm-queries/users/transaction/:id", hndlr.UpdateUserWithTransaction)
	group.Put("/orm-queries/users/selected-columns/:id", hndlr.UpdateUserWithSelectedColumns)
	group.Put("/orm-queries/users/returning/:id", hndlr.UpdateUserReturningResult)
	group.Put("/native-queries/users/sql/:id", hndlr.UpdateUserWithNativeSQL)
	group.Put("/native-queries/users/mongodb/:id", hndlr.UpdateUserWithNativeMongoQuery)
	group.Put("/native-queries/users/scylla/:id", hndlr.UpdateUserWithNativeCQL)
	group.Put("/orm-queries/users/:id", hndlr.UpdateUser)

	// PATCH changes only the isActive property. Example body: {"isActive":false}.
	group.Patch("/orm-queries/users/:id", hndlr.PatchUserActiveStatus)

	// Delete users. Keep fixed paths before the /:id parameter route.
	group.Delete("/orm-queries/users/transaction/:id", hndlr.DeleteUserWithTransaction)
	group.Delete("/native-queries/users/sql/:id", hndlr.DeleteUserWithNativeSQL)
	group.Delete("/native-queries/users/mongodb/:id", hndlr.DeleteUserWithNativeMongoQuery)
	group.Delete("/native-queries/users/scylla/:id", hndlr.DeleteUserWithNativeCQL)
	group.Delete("/orm-queries/users/:id", hndlr.DeleteUser)
}

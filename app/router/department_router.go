package router

import (
	"net/http"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest/server"
)

// DepartmentRouter registers department CRUD examples under /examples.
func DepartmentRouter(app *gorest.App) {
	mdl := model.NewDepartmentModel(app)
	lib := library.NewDepartmentLibrary(mdl)
	hndlr := handler.NewDepartmentHandler(app, lib)
	RegisterDepartmentRoutes(app.Group("/examples"), hndlr)
}

// RegisterDepartmentRoutes registers department endpoints on a route group.
func RegisterDepartmentRoutes(group *server.Server, hndlr *handler.DepartmentHandler) {
	// Read the collection or one department.
	group.Get("/departments", hndlr.GetDepartments)
	group.Get("/departments/:id", hndlr.GetDepartmentByID)

	// HEAD reuses the GET handlers; the response has no body.
	group.Head("/departments", hndlr.GetDepartments)
	group.Head("/departments/:id", hndlr.GetDepartmentByID)

	// OPTIONS advertises methods available on the collection and one department.
	group.Options("/departments", func(ctx *gorest.Context) error {
		ctx.Set("Allow", "GET, HEAD, POST, OPTIONS")
		return ctx.SendStatus(http.StatusNoContent)
	})
	group.Options("/departments/:id", func(ctx *gorest.Context) error {
		ctx.Set("Allow", "GET, HEAD, PUT, PATCH, DELETE, OPTIONS")
		return ctx.SendStatus(http.StatusNoContent)
	})

	// Department create, update, and delete examples call BeginTx. GOREST's
	// BeginTx/Commit/Rollback calls are no-ops for MongoDB and ScyllaDB.
	// Create a department: POST /examples/departments
	group.Post("/departments", hndlr.CreateDepartment)

	// Replace department fields: PUT /examples/departments/:id
	group.Put("/departments/:id", hndlr.UpdateDepartment)

	// Partially update only fields present in the JSON body, for example:
	// PATCH /examples/departments/:id with {"isActive":false}.
	group.Patch("/departments/:id", hndlr.PatchDepartment)

	// Delete one department: DELETE /examples/departments/:id
	group.Delete("/departments/:id", hndlr.DeleteDepartment)
}

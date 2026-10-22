package router

import (
	"net/http"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest/server"
)

// OrderRouter registers order CRUD examples under /examples.
func OrderRouter(app *gorest.App) {
	mdl := model.NewOrderModel(app)
	lib := library.NewOrderLibrary(mdl)
	hndlr := handler.NewOrderHandler(app, lib)
	RegisterOrderRoutes(app.Group("/examples"), hndlr)
}

// RegisterOrderRoutes registers order endpoints on a route group.
func RegisterOrderRoutes(group *server.Server, hndlr *handler.OrderHandler) {
	// Read the collection or one order.
	group.Get("/orders", hndlr.GetOrders)
	group.Get("/users/:userId/orders", hndlr.GetOrdersByUserID)
	group.Get("/orders/:id", hndlr.GetOrderByID)

	// HEAD reuses the GET handlers; the response has no body.
	group.Head("/orders", hndlr.GetOrders)
	group.Head("/users/:userId/orders", hndlr.GetOrdersByUserID)
	group.Head("/orders/:id", hndlr.GetOrderByID)

	// OPTIONS advertises methods available on the collection and one order.
	group.Options("/orders", func(ctx *gorest.Context) error {
		ctx.Set("Allow", "GET, HEAD, POST, OPTIONS")
		return ctx.SendStatus(http.StatusNoContent)
	})
	group.Options("/users/:userId/orders", func(ctx *gorest.Context) error {
		ctx.Set("Allow", "GET, HEAD, OPTIONS")
		return ctx.SendStatus(http.StatusNoContent)
	})
	group.Options("/orders/:id", func(ctx *gorest.Context) error {
		ctx.Set("Allow", "GET, HEAD, PUT, PATCH, DELETE, OPTIONS")
		return ctx.SendStatus(http.StatusNoContent)
	})

	// Create, update, and delete use BeginTx. MongoDB and ScyllaDB do not provide
	// transaction atomicity through GOREST's BeginTx/Commit/Rollback methods.
	// Create an order: POST /examples/orders
	group.Post("/orders", hndlr.CreateOrder)

	// Replace order fields: PUT /examples/orders/:id
	group.Put("/orders/:id", hndlr.UpdateOrder)

	// Partially update only fields present in the JSON body, for example:
	// PATCH /examples/orders/:id with {"isActive":false}.
	group.Patch("/orders/:id", hndlr.PatchOrder)

	// Delete one order: DELETE /examples/orders/:id
	group.Delete("/orders/:id", hndlr.DeleteOrder)
}

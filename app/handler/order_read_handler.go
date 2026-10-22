package handler

import (
	"github.com/elsyahtech/gorest"
	resview "github.com/elsyahtech/gorest-api-starter/app/view"
)

// GetOrders handles a request to list orders.
func (handler *OrderHandler) GetOrders(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.OrderLibrary.GetOrders()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildOrderResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetOrderByID handles a request to retrieve one order by ID.
func (handler *OrderHandler) GetOrderByID(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	orderID := ctx.Params("id")

	data, message, httpCode, err := handler.OrderLibrary.GetOrderByID(orderID)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildOrderResponse(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetOrdersByUserID handles a request to list orders for one user.
func (handler *OrderHandler) GetOrdersByUserID(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.OrderLibrary.GetOrdersByUserID(ctx.Params("userId"))
	if err != nil {
		response.WithMessage(message, err, httpCode)
		handler.App.Log(gorest.Map{"error": err}).Error(message)
		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildOrderResponseArray(data)
	response.WithMessage(message, nil, httpCode)
	return response.SendJSON(ctx, response, nil)
}

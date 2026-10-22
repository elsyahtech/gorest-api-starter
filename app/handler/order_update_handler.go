package handler

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

// UpdateOrder handles a request to replace order data.
func (handler *OrderHandler) UpdateOrder(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	var request entity.Order

	message, httpCode, err := helper.BodyRequest(ctx, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	orderID := ctx.Params("id")

	message, httpCode, err = handler.OrderLibrary.UpdateOrder(orderID, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// PatchOrder handles a partial update to order data.
func (handler *OrderHandler) PatchOrder(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	var request entity.OrderPatch

	message, httpCode, err := helper.BodyRequest(ctx, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)
		handler.App.Log(gorest.Map{"error": err}).Error(message)
		return response.SendJSON(ctx, response, err)
	}

	message, httpCode, err = handler.OrderLibrary.PatchOrder(ctx.Params("id"), &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)
		handler.App.Log(gorest.Map{"error": err}).Error(message)
		return response.SendJSON(ctx, response, err)
	}

	response.WithMessage(message, nil, httpCode)
	return response.SendJSON(ctx, response, nil)
}

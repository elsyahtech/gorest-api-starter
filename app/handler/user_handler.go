package handler

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	resview "github.com/elsyahtech/gorest-api-starter/app/view"
	"github.com/elsyahtech/gorest/helper"
)

func (hndlr *UserHandler) Login(ctx *gorest.Context) error {
	response := hndlr.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	var request entity.User

	message, httpCode, err := helper.BodyRequest(ctx, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		hndlr.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err) //nolint:wrapcheck // reason: wrapcheck is disabled for this line to avoid unnecessary error wrapping.
	}

	message, httpCode, err = hndlr.UserLibrary.Login(ctx, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		hndlr.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err) //nolint:wrapcheck // reason: wrapcheck is disabled for this line to avoid unnecessary error wrapping.
	}

	response.WithMessage("Login successful", nil, httpCode)

	return response.SendJSON(ctx, response, nil) //nolint:wrapcheck // reason: wrapcheck is disabled for this line to avoid unnecessary error wrapping.
}

// CreateUser handles a request to create a user through the ORM.
func (handler *UserHandler) CreateUser(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	var request entity.User

	message, httpCode, err := helper.BodyRequest(ctx, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err) //nolint:wrapcheck // reason: wrapcheck is disabled for this line to avoid unnecessary error wrapping.
	}

	message, httpCode, err = handler.UserLibrary.CreateUser(&request)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err) //nolint:wrapcheck // reason: wrapcheck is disabled for this line to avoid unnecessary error wrapping.
	}

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil) //nolint:wrapcheck // reason: wrapcheck is disabled for this line to avoid unnecessary error wrapping.
}

func (handler *UserHandler) GetUsers(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsers()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err) //nolint:wrapcheck // reason: wrapcheck is disabled for this line to avoid unnecessary error wrapping.
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil) //nolint:wrapcheck // reason: wrapcheck is disabled for this line to avoid unnecessary error wrapping.
}

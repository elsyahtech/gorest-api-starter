package handler

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

// Login handles user login requests. It validates the provided credentials and generates an
// authentication token upon successful login.
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

		return response.SendJSON(ctx, response, err)
	}

	message, httpCode, err = hndlr.UserLibrary.Login(ctx, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		hndlr.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.WithMessage("Login successful", nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

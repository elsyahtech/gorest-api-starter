package handler

import (
	"github.com/elsyahtech/gorest"
)

// DeleteUserWithTransaction demonstrates deleting a user within the configured transaction flow.
func (handler *UserHandler) DeleteUserWithTransaction(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	userID := ctx.Params("id")

	message, httpCode, err := handler.UserLibrary.DeleteUserWithTransaction(userID)
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

// DeleteUserWithNativeSQL demonstrates deleting a user with a native SQL statement.
func (handler *UserHandler) DeleteUserWithNativeSQL(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	userID := ctx.Params("id")

	message, httpCode, err := handler.UserLibrary.DeleteUserWithNativeSQL(userID)
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

// DeleteUserWithNativeMongoQuery demonstrates deleting a user with a native MongoDB query.
func (handler *UserHandler) DeleteUserWithNativeMongoQuery(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	userID := ctx.Params("id")

	message, httpCode, err := handler.UserLibrary.DeleteUserWithNativeMongoQuery(userID)
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

// DeleteUserWithNativeCQL demonstrates deleting a user with a native ScyllaDB CQL statement.
func (handler *UserHandler) DeleteUserWithNativeCQL(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	userID := ctx.Params("id")

	message, httpCode, err := handler.UserLibrary.DeleteUserWithNativeCQL(userID)
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

// DeleteUser handles a request to delete a user through the ORM.
func (handler *UserHandler) DeleteUser(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	userID := ctx.Params("id")

	message, httpCode, err := handler.UserLibrary.DeleteUser(userID)
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

// DeleteUserInRedis handles a request to delete a user in Redis.
func (handler *UserHandler) DeleteUserInRedis(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	userID := ctx.Params("id")

	message, httpCode, err := handler.UserLibrary.DeleteUserInRedis(userID)
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

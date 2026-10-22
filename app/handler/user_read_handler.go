package handler

import (
	"github.com/elsyahtech/gorest"
	resview "github.com/elsyahtech/gorest-api-starter/app/view"
	"github.com/elsyahtech/gorest/helper"
)

// GetUsers handles a request to list users.
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

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersWithAllColumns demonstrates reading users with all columns selected.
func (handler *UserHandler) GetUsersWithAllColumns(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersWithAllColumns()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersWithSelectedColumns demonstrates reading users with a column projection.
func (handler *UserHandler) GetUsersWithSelectedColumns(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersWithSelectedColumns()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersOrderedByID demonstrates listing users in descending ID order.
func (handler *UserHandler) GetUsersOrderedByID(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersOrderedByID()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUserCountsByDepartment demonstrates grouping users and counting by department.
func (handler *UserHandler) GetUserCountsByDepartment(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUserCountsByDepartment()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildGroupByResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersWithDistinctDepartmentIDs demonstrates querying distinct department IDs from users.
func (handler *UserHandler) GetUsersWithDistinctDepartmentIDs(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersWithDistinctDepartmentIDs()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetDepartmentsByUserCount demonstrates filtering departments with more than two users using HAVING.
func (handler *UserHandler) GetDepartmentsByUserCount(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetDepartmentsByUserCount()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildGroupByResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersWithDepartment demonstrates loading users with a SQL join.
func (handler *UserHandler) GetUsersWithDepartment(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersWithDepartment()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.ToJoinResponses(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersWithPreloadedDepartment demonstrates loading a user relation with Preload.
func (handler *UserHandler) GetUsersWithPreloadedDepartment(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersWithPreloadedDepartment()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersWithNativeSQL demonstrates reading users with a native SQL query.
func (handler *UserHandler) GetUsersWithNativeSQL(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersWithNativeSQL()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersWithNativeMongoQuery demonstrates querying active users with a native MongoDB query.
func (handler *UserHandler) GetUsersWithNativeMongoQuery(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersWithNativeMongoQuery()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersWithNativeCQL demonstrates reading users with a native ScyllaDB CQL query.
func (handler *UserHandler) GetUsersWithNativeCQL(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.UserLibrary.GetUsersWithNativeCQL()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersByLimit demonstrates limiting the number of users returned.
func (handler *UserHandler) GetUsersByLimit(ctx *gorest.Context) error {
	var gorestMap gorest.Map

	response := handler.JSON.NewResponse(gorestMap)

	limitParam := ctx.Query("limit")

	limit, message, httpCode, err := helper.ConvertStringToInt("limit", limitParam)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	data, message, httpCode, err := handler.UserLibrary.GetUsersByLimit(limit)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	gorestMap = gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
		"Total":     len(*data),
	}

	response = handler.JSON.NewResponse(gorestMap)

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUsersByLimitOffset demonstrates paginating users with limit and offset.
func (handler *UserHandler) GetUsersByLimitOffset(ctx *gorest.Context) error {
	var gorestMap gorest.Map

	response := handler.JSON.NewResponse(gorestMap)

	limitParam := ctx.Query("limit")

	limit, message, httpCode, err := helper.ConvertStringToInt("limit", limitParam)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	offsetParam := ctx.Query("offset")

	offset, message, httpCode, err := helper.ConvertStringToInt("offset", offsetParam)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	data, message, httpCode, err := handler.UserLibrary.GetUsersByLimitOffset(limit, offset)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	gorestMap = gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
		"Total":     len(*data),
	}

	response = handler.JSON.NewResponse(gorestMap)

	response.Data = resview.BuildUserResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetUserByID handles a request to retrieve one user by ID.
func (handler *UserHandler) GetUserByID(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	userID := ctx.Params("id")

	data, message, httpCode, err := handler.UserLibrary.GetUserByID(userID)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildUserResponse(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

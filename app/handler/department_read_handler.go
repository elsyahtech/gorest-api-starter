package handler

import (
	"github.com/elsyahtech/gorest"
	resview "github.com/elsyahtech/gorest-api-starter/app/view"
)

// GetDepartments handles a request to list departments.
func (handler *DepartmentHandler) GetDepartments(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	data, message, httpCode, err := handler.DepartmentLibrary.GetDepartments()
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildDepartmentResponseArray(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

// GetDepartmentByID handles a request to retrieve one department by ID.
func (handler *DepartmentHandler) GetDepartmentByID(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	departmentID := ctx.Params("id")

	data, message, httpCode, err := handler.DepartmentLibrary.GetDepartmentByID(departmentID)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	response.Data = resview.BuildDepartmentResponse(data)

	response.WithMessage(message, nil, httpCode)

	return response.SendJSON(ctx, response, nil)
}

package handler

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

// UpdateDepartment handles a request to replace department data.
func (handler *DepartmentHandler) UpdateDepartment(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	var request entity.Department

	message, httpCode, err := helper.BodyRequest(ctx, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)

		handler.App.Log(gorest.Map{
			"error": err,
		}).Error(message)

		return response.SendJSON(ctx, response, err)
	}

	departmentID := ctx.Params("id")

	message, httpCode, err = handler.DepartmentLibrary.UpdateDepartment(departmentID, &request)
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

// PatchDepartment handles a partial update to department data.
func (handler *DepartmentHandler) PatchDepartment(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	var request entity.DepartmentPatch

	message, httpCode, err := helper.BodyRequest(ctx, &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)
		handler.App.Log(gorest.Map{"error": err}).Error(message)
		return response.SendJSON(ctx, response, err)
	}

	message, httpCode, err = handler.DepartmentLibrary.PatchDepartment(ctx.Params("id"), &request)
	if err != nil {
		response.WithMessage(message, err, httpCode)
		handler.App.Log(gorest.Map{"error": err}).Error(message)
		return response.SendJSON(ctx, response, err)
	}

	response.WithMessage(message, nil, httpCode)
	return response.SendJSON(ctx, response, nil)
}

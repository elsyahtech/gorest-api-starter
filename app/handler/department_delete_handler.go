package handler

import "github.com/elsyahtech/gorest"

// DeleteDepartment handles a request to delete a department.
func (handler *DepartmentHandler) DeleteDepartment(ctx *gorest.Context) error {
	response := handler.JSON.NewResponse(gorest.Map{
		"developer": "Elmansyah",
		"contact":   "https://www.linkedin.com/in/elmansyah",
		"version":   "v1",
	})

	departmentID := ctx.Params("id")

	message, httpCode, err := handler.DepartmentLibrary.DeleteDepartment(departmentID)
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

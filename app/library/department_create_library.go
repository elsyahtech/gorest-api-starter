package library

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// CreateDepartment validates and creates a department.
func (lib *DepartmentLibrary) CreateDepartment(req *entity.Department) (string, int, error) {
	if req == nil {
		return "Department payload must not be empty", http.StatusBadRequest,
			errors.New("add department: request is empty")
	}

	message, httpCode, err := lib.DepartmentModel.CreateDepartment(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

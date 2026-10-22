package library

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// GetDepartments returns all departments.
func (lib *DepartmentLibrary) GetDepartments() (*entity.DepartmentArray, string, int, error) {
	result, message, httpCode, err := lib.DepartmentModel.GetDepartments()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetDepartmentByID returns the department with the given ID.
func (lib *DepartmentLibrary) GetDepartmentByID(departmentID string) (*entity.Department, string, int, error) {
	if departmentID == "" {
		return nil, "Department ID param must not be empty", http.StatusBadRequest,
			errors.New("get department by ID: department ID is empty")
	}

	result, message, httpCode, err := lib.DepartmentModel.GetDepartmentByID(departmentID)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

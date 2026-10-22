package library

import (
	"errors"
	"fmt"
	"net/http"
)

// DeleteDepartment removes the department with the given ID.
func (lib *DepartmentLibrary) DeleteDepartment(departmentID string) (string, int, error) {
	if departmentID == "" {
		return "Department ID param must not be empty", http.StatusBadRequest,
			errors.New("delete department: department ID is empty")
	}

	message, httpCode, err := lib.DepartmentModel.DeleteDepartment(departmentID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

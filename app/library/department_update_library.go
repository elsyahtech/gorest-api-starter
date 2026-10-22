package library

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// UpdateDepartment validates and replaces a department.
func (lib *DepartmentLibrary) UpdateDepartment(departmentID string, req *entity.Department) (string, int, error) {
	if departmentID == "" {
		return "Department ID param must not be empty", http.StatusBadRequest,
			errors.New("update department: department ID is empty")
	}

	if req == nil {
		return "Department payload must not be empty", http.StatusBadRequest,
			errors.New("update department: payload is empty")
	}

	message, httpCode, err := lib.DepartmentModel.UpdateDepartment(departmentID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// PatchDepartment applies a partial update to a department.
func (lib *DepartmentLibrary) PatchDepartment(departmentID string, req *entity.DepartmentPatch) (string, int, error) {
	if departmentID == "" {
		return "Department ID param must not be empty", http.StatusBadRequest,
			errors.New("patch department: department ID is empty")
	}

	if req == nil || (req.Name == nil && req.Description == nil && req.IsActive == nil) {
		return "Provide at least one department field to patch", http.StatusBadRequest,
			errors.New("patch department: no fields provided")
	}

	message, httpCode, err := lib.DepartmentModel.PatchDepartment(departmentID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("patch department: %w", err)
	}

	return message, httpCode, nil
}

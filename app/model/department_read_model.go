package model

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// GetDepartments returns all departments.
func (model *DepartmentModel) GetDepartments() (*entity.DepartmentArray, string, int, error) {
	database := model.App.Database()

	var departments entity.DepartmentArray

	query := database.Table("departments").Find(&departments)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &departments, query.Message, query.HTTPCode, nil
}

// GetDepartmentByID returns the department with the given ID.
func (model *DepartmentModel) GetDepartmentByID(departmentID string) (*entity.Department, string, int, error) {
	if departmentID == "" {
		return nil, "Department ID param must not be empty", http.StatusBadRequest,
			errors.New("get department by ID: department ID is empty")
	}

	database := model.App.Database()

	var department entity.Department

	query := database.Table("departments").Where("id = ?", departmentID).Find(&department)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &department, "", query.HTTPCode, nil
}

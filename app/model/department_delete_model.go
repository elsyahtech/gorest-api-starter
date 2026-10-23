package model

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// DeleteDepartment removes the department with the given ID in a transaction.
func (model *DepartmentModel) DeleteDepartment(departmentID string) (string, int, error) {
	if departmentID == "" {
		return "Department ID must not be empty", http.StatusBadRequest,
			errors.New("delete department: department ID is empty")
	}

	database := model.App.Database()

	transaction, message, httpCode, err := database.BeginTx()
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	defer func() {
		if message, _, err := transaction.Rollback(); err != nil {
			model.App.Log(gorest.Map{"error": err}).Error(message)
		}
	}()

	var department entity.Department

	query := transaction.Table("departments").Where("id = ?", departmentID).Delete(&department)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

package model

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// UpdateDepartment updates the department with the given ID in a transaction.
func (model *DepartmentModel) UpdateDepartment(departmentID string, req *entity.Department) (string, int, error) {
	if departmentID == "" {
		return "Department ID param must not be empty", http.StatusBadRequest,
			errors.New("update department: department ID is empty")
	}

	if req == nil {
		return "Department payload must not be empty", http.StatusBadRequest,
			errors.New("update department: payload is empty")
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

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := transaction.Table("departments").Where("id = ?", departmentID).Update(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

// PatchDepartment updates only the department fields present in the patch request.
func (model *DepartmentModel) PatchDepartment(departmentID string, req *entity.DepartmentPatch) (string, int, error) {
	if departmentID == "" {
		return "Department ID param must not be empty", http.StatusBadRequest,
			errors.New("patch department: department ID is empty")
	}

	if req == nil {
		return "Department patch must not be empty", http.StatusBadRequest,
			errors.New("patch department: payload is empty")
	}

	updated := entity.Department{}
	columns := make([]string, 0, 4)

	if req.Name != nil {
		updated.Name = *req.Name

		columns = append(columns, "name")
	}

	if req.Description != nil {
		updated.Description = *req.Description

		columns = append(columns, "description")
	}

	if req.IsActive != nil {
		updated.IsActive = *req.IsActive

		columns = append(columns, "is_active")
	}

	if len(columns) == 0 {
		return "Provide at least one department field to patch", http.StatusBadRequest,
			errors.New("patch department: no fields provided")
	}

	now := time.Now()
	updated.UpdatedAt = &now

	columns = append(columns, "updated_at")

	query := model.App.Database().
		Table("departments").
		Select(columns...).
		Where("id = ?", departmentID).
		Update(&updated)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("patch department: %w", query.Error)
	}

	return "Department updated successfully", query.HTTPCode, nil
}

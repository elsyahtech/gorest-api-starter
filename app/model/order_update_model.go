package model

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// UpdateOrder updates the order with the given ID in a transaction.
func (model *OrderModel) UpdateOrder(orderID string, req *entity.Order) (string, int, error) {
	if orderID == "" {
		return "Order ID param must not be empty", http.StatusBadRequest,
			errors.New("update order: order ID is empty")
	}

	if req == nil {
		return "Order payload must not be empty", http.StatusBadRequest,
			errors.New("update order: payload is empty")
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

	req.UpdatedAt = &timeNow

	query := transaction.Table("orders").Where("id = ?", orderID).Update(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

// PatchOrder updates only the order fields present in the patch request.
func (model *OrderModel) PatchOrder(orderID string, req *entity.OrderPatch) (string, int, error) {
	if orderID == "" {
		return "Order ID param must not be empty", http.StatusBadRequest,
			errors.New("patch order: order ID is empty")
	}

	if req == nil {
		return "Order patch must not be empty", http.StatusBadRequest,
			errors.New("patch order: payload is empty")
	}

	updated := entity.Order{}
	columns := make([]string, 0, 4)

	if req.OrderNumber != nil {
		updated.OrderNumber = *req.OrderNumber

		columns = append(columns, "order_number")
	}

	if req.Status != nil {
		updated.Status = *req.Status

		columns = append(columns, "status")
	}

	if req.TotalAmount != nil {
		updated.TotalAmount = *req.TotalAmount

		columns = append(columns, "total_amount")
	}

	if req.IsActive != nil {
		updated.IsActive = *req.IsActive

		columns = append(columns, "is_active")
	}

	if len(columns) == 0 {
		return "Provide at least one order field to patch", http.StatusBadRequest,
			errors.New("patch order: no fields provided")
	}

	now := time.Now()

	updated.UpdatedAt = &now

	columns = append(columns, "updated_at")

	query := model.App.Database().
		Table("orders").
		Select(columns...).
		Where("id = ?", orderID).
		Update(&updated)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("patch order: %w", query.Error)
	}

	return "Order updated successfully", query.HTTPCode, nil
}

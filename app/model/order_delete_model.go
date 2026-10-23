package model

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// DeleteOrder removes the order with the given ID in a transaction.
func (model *OrderModel) DeleteOrder(orderID string) (string, int, error) {
	if orderID == "" {
		return "Order ID must not be empty", http.StatusBadRequest,
			errors.New("delete order: order ID is empty")
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

	var order entity.Order

	query := transaction.Table("orders").Where("id = ?", orderID).Delete(&order)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

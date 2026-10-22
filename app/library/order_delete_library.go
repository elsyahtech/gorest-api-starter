package library

import (
	"errors"
	"fmt"
	"net/http"
)

// DeleteOrder removes the order with the given ID.
func (lib *OrderLibrary) DeleteOrder(orderID string) (string, int, error) {
	if orderID == "" {
		return "Order ID param must not be empty", http.StatusBadRequest,
			errors.New("delete order: order ID is empty")
	}

	message, httpCode, err := lib.OrderModel.DeleteOrder(orderID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

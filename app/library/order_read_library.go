package library

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// GetOrders returns all orders.
func (lib *OrderLibrary) GetOrders() (*entity.OrderArray, string, int, error) {
	result, message, httpCode, err := lib.OrderModel.GetOrders()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetOrderByID returns the order with the given ID.
func (lib *OrderLibrary) GetOrderByID(orderID string) (*entity.Order, string, int, error) {
	if orderID == "" {
		return nil, "Order ID param must not be empty", http.StatusBadRequest,
			errors.New("get order by ID: order ID is empty")
	}

	result, message, httpCode, err := lib.OrderModel.GetOrderByID(orderID)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetOrdersByUserID returns orders that belong to one user.
func (lib *OrderLibrary) GetOrdersByUserID(userID string) (*entity.OrderArray, string, int, error) {
	if userID == "" {
		return nil, "User ID param must not be empty", http.StatusBadRequest,
			errors.New("get orders by user ID: user ID is empty")
	}

	result, message, httpCode, err := lib.OrderModel.GetOrdersByUserID(userID)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("get orders by user ID: %w", err)
	}

	return result, message, httpCode, nil
}

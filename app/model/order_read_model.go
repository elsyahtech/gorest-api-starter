package model

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// GetOrders returns all orders.
func (model *OrderModel) GetOrders() (*entity.OrderArray, string, int, error) {
	database := model.App.Database()

	var orders entity.OrderArray

	query := database.Table("orders").Find(&orders)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &orders, query.Message, query.HTTPCode, nil
}

// GetOrderByID returns the order with the given ID.
func (model *OrderModel) GetOrderByID(orderID string) (*entity.Order, string, int, error) {
	if orderID == "" {
		return nil, "Order ID param must not be empty", http.StatusBadRequest,
			errors.New("get order by ID: order ID is empty")
	}

	database := model.App.Database()

	var order entity.Order

	query := database.Table("orders").Where("id = ?", orderID).Find(&order)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &order, "", query.HTTPCode, nil
}

// GetOrdersByUserID returns orders belonging to the given user.
func (model *OrderModel) GetOrdersByUserID(userID string) (*entity.OrderArray, string, int, error) {
	if userID == "" {
		return nil, "User ID param must not be empty", http.StatusBadRequest,
			errors.New("get orders by user ID: user ID is empty")
	}

	database := model.App.Database()

	var orders entity.OrderArray

	query := database.Table("orders").Where("user_id = ?", userID).Find(&orders)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &orders, query.Message, query.HTTPCode, nil
}

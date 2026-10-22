package library

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// UpdateOrder validates and replaces an order.
func (lib *OrderLibrary) UpdateOrder(orderID string, req *entity.Order) (string, int, error) {
	if orderID == "" {
		return "Order ID param must not be empty", http.StatusBadRequest,
			errors.New("update order: order ID is empty")
	}

	if req == nil {
		return "Order payload must not be empty", http.StatusBadRequest,
			errors.New("update order: payload is empty")
	}

	message, httpCode, err := lib.OrderModel.UpdateOrder(orderID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// PatchOrder applies a partial update to an order.
func (lib *OrderLibrary) PatchOrder(orderID string, req *entity.OrderPatch) (string, int, error) {
	if orderID == "" {
		return "Order ID param must not be empty", http.StatusBadRequest,
			errors.New("patch order: order ID is empty")
	}

	if req == nil || (req.OrderNumber == nil && req.Status == nil && req.TotalAmount == nil && req.IsActive == nil) {
		return "Provide at least one order field to patch", http.StatusBadRequest,
			errors.New("patch order: no fields provided")
	}

	message, httpCode, err := lib.OrderModel.PatchOrder(orderID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("patch order: %w", err)
	}

	return message, httpCode, nil
}

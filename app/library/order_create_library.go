package library

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// CreateOrder validates and creates an order.
func (lib *OrderLibrary) CreateOrder(req *entity.Order) (string, int, error) {
	if req == nil {
		return "Order payload must not be empty", http.StatusBadRequest,
			errors.New("add order: request is empty")
	}

	message, httpCode, err := lib.OrderModel.CreateOrder(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

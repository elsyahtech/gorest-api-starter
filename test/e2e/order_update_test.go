package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestUpdateOrderE2E(t *testing.T) {
	testServer, userID := newOrderTestServer(t)
	orderID := "order-update-e2e-01"
	initial := entity.Order{
		ID:          orderID,
		UserID:      userID,
		OrderNumber: "ORD-E2E-UPDATE-OLD",
		Status:      "PENDING",
		TotalAmount: 45.00,
		IsActive:    true,
	}
	createBody, err := json.Marshal(initial)
	if err != nil {
		t.Fatalf("encode create order request: %v", err)
	}

	createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/orders", createBody)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orders status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	updated := entity.Order{
		ID:          orderID,
		UserID:      userID,
		OrderNumber: "ORD-E2E-UPDATE-NEW",
		Status:      "FULFILLED",
		TotalAmount: 145.75,
		IsActive:    false,
	}
	updateBody, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode update order request: %v", err)
	}

	updateResponse := sendRequest(
		t,
		testServer,
		http.MethodPut,
		"/examples/orders/"+orderID,
		updateBody,
	)
	if updateResponse.StatusCode < http.StatusOK || updateResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PUT /examples/orders/%s status = %d, want 2xx; body: %s",
			orderID, updateResponse.StatusCode, updateResponse.Body)
	}

	var updateResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(updateResponse.Body, &updateResult); err != nil {
		t.Fatalf("decode update order response: %v", err)
	}
	if !updateResult.Success {
		t.Fatal("update order response success = false, want true")
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orders/"+orderID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orders/%s status = %d, want %d; body: %s",
			orderID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}

	var getResult struct {
		Data entity.OrderResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &getResult); err != nil {
		t.Fatalf("decode order response after update: %v", err)
	}
	if getResult.Data.ID != updated.ID ||
		getResult.Data.UserID != updated.UserID ||
		getResult.Data.OrderNumber != updated.OrderNumber ||
		getResult.Data.Status != updated.Status ||
		getResult.Data.TotalAmount != updated.TotalAmount ||
		getResult.Data.IsActive != updated.IsActive {
		t.Fatalf("order after update = %+v, want ID %q, userID %q, orderNumber %q, status %q, totalAmount %.2f, isActive %t",
			getResult.Data, updated.ID, updated.UserID, updated.OrderNumber,
			updated.Status, updated.TotalAmount, updated.IsActive)
	}
	if getResult.Data.UpdatedAt == nil || *getResult.Data.UpdatedAt == "" {
		t.Fatal("order updatedAt is empty, want a formatted update timestamp")
	}
}

package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestCreateOrderE2E(t *testing.T) {
	testServer, userID := newOrderTestServer(t)
	expected := entity.Order{
		ID:          "order-create-e2e-01",
		UserID:      userID,
		OrderNumber: "ORD-E2E-0001",
		Status:      "PENDING",
		TotalAmount: 125.50,
		IsActive:    true,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("encode create order request: %v", err)
	}

	createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/orders", body)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orders status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode create order response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("create order response success = false, want true")
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orders/"+expected.ID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orders/%s status = %d, want %d; body: %s",
			expected.ID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}

	var getResult struct {
		Data entity.OrderResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &getResult); err != nil {
		t.Fatalf("decode order response: %v", err)
	}
	if getResult.Data.ID != expected.ID ||
		getResult.Data.UserID != expected.UserID ||
		getResult.Data.OrderNumber != expected.OrderNumber ||
		getResult.Data.Status != expected.Status ||
		getResult.Data.TotalAmount != expected.TotalAmount ||
		getResult.Data.IsActive != expected.IsActive {
		t.Fatalf("order response = %+v, want ID %q, userID %q, orderNumber %q, status %q, totalAmount %.2f, isActive %t",
			getResult.Data, expected.ID, expected.UserID, expected.OrderNumber, expected.Status,
			expected.TotalAmount, expected.IsActive)
	}
	if getResult.Data.CreatedAt == "" {
		t.Fatal("order createdAt is empty, want a formatted creation timestamp")
	}
}

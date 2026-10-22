package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestDeleteOrderE2E(t *testing.T) {
	testServer, userID := newOrderTestServer(t)
	order := entity.Order{
		ID:          "order-delete-e2e-01",
		UserID:      userID,
		OrderNumber: "ORD-E2E-DELETE-0001",
		Status:      "PENDING",
		TotalAmount: 75.25,
		IsActive:    true,
	}
	body, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("encode create order request: %v", err)
	}

	createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/orders", body)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orders status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	deleteResponse := sendRequest(
		t,
		testServer,
		http.MethodDelete,
		"/examples/orders/"+order.ID,
		nil,
	)
	if deleteResponse.StatusCode < http.StatusOK || deleteResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("DELETE /examples/orders/%s status = %d, want 2xx; body: %s",
			order.ID, deleteResponse.StatusCode, deleteResponse.Body)
	}

	var deleteResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(deleteResponse.Body, &deleteResult); err != nil {
		t.Fatalf("decode delete order response: %v", err)
	}
	if !deleteResult.Success {
		t.Fatal("delete order response success = false, want true")
	}

	listResponse := sendRequest(t, testServer, http.MethodGet, "/examples/orders", nil)
	if listResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orders status = %d, want %d; body: %s",
			listResponse.StatusCode, http.StatusOK, listResponse.Body)
	}

	var listResult struct {
		Data entity.OrderResponseArray `json:"data"`
	}
	if err := json.Unmarshal(listResponse.Body, &listResult); err != nil {
		t.Fatalf("decode order list response after delete: %v", err)
	}
	if len(listResult.Data) != 0 {
		t.Fatalf("order count after delete = %d, want 0", len(listResult.Data))
	}
}

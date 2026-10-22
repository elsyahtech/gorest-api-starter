package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestPatchOrderE2E(t *testing.T) {
	testServer, userID := newOrderTestServer(t)
	orderID := "order-patch-e2e-01"
	initial := entity.Order{
		ID:          orderID,
		UserID:      userID,
		OrderNumber: "ORD-E2E-PATCH-0001",
		Status:      "PENDING",
		TotalAmount: 225.00,
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

	orderNumber := "ORD-E2E-PATCH-UPDATED"
	totalAmount := float64(0)
	isActive := false
	patchBody, err := json.Marshal(entity.OrderPatch{
		OrderNumber: &orderNumber,
		TotalAmount: &totalAmount,
		IsActive:    &isActive,
	})
	if err != nil {
		t.Fatalf("encode patch order request: %v", err)
	}

	patchResponse := sendRequest(
		t,
		testServer,
		http.MethodPatch,
		"/examples/orders/"+orderID,
		patchBody,
	)
	if patchResponse.StatusCode < http.StatusOK || patchResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PATCH /examples/orders/%s status = %d, want 2xx; body: %s",
			orderID, patchResponse.StatusCode, patchResponse.Body)
	}

	var patchResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(patchResponse.Body, &patchResult); err != nil {
		t.Fatalf("decode patch order response: %v", err)
	}
	if !patchResult.Success {
		t.Fatal("patch order response success = false, want true")
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
		t.Fatalf("decode order response after patch: %v", err)
	}
	if getResult.Data.ID != initial.ID ||
		getResult.Data.UserID != initial.UserID ||
		getResult.Data.OrderNumber != orderNumber ||
		getResult.Data.Status != initial.Status ||
		getResult.Data.TotalAmount != totalAmount ||
		getResult.Data.IsActive {
		t.Fatalf("order after patch = %+v, want updated orderNumber, totalAmount=0, isActive=false and other fields unchanged",
			getResult.Data)
	}
	if getResult.Data.UpdatedAt == nil || *getResult.Data.UpdatedAt == "" {
		t.Fatal("order updatedAt is empty, want a formatted update timestamp")
	}

	emptyPatchResponse := sendRequest(
		t,
		testServer,
		http.MethodPatch,
		"/examples/orders/"+orderID,
		[]byte(`{}`),
	)
	if emptyPatchResponse.StatusCode != http.StatusBadRequest {
		t.Fatalf("PATCH with empty object status = %d, want %d; body: %s",
			emptyPatchResponse.StatusCode, http.StatusBadRequest, emptyPatchResponse.Body)
	}
}

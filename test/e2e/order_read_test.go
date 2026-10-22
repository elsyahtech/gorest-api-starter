package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetOrdersE2E(t *testing.T) {
	testServer, userID := newOrderTestServer(t)

	emptyResponse := sendRequest(t, testServer, http.MethodGet, "/examples/orders", nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orders status = %d, want %d; body: %s",
			emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}

	var emptyResult struct {
		Data entity.OrderResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty order list response: %v", err)
	}
	if emptyResult.Data == nil {
		t.Fatal("empty order list data is null, want an empty array")
	}
	if len(emptyResult.Data) != 0 {
		t.Fatalf("order count in empty list = %d, want 0", len(emptyResult.Data))
	}

	expectedOrders := []entity.Order{
		{
			ID:          "order-read-e2e-01",
			UserID:      userID,
			OrderNumber: "ORD-E2E-READ-0001",
			Status:      "PENDING",
			TotalAmount: 49.99,
			IsActive:    true,
		},
		{
			ID:          "order-read-e2e-02",
			UserID:      userID,
			OrderNumber: "ORD-E2E-READ-0002",
			Status:      "FULFILLED",
			TotalAmount: 175.25,
			IsActive:    false,
		},
	}
	for _, order := range expectedOrders {
		body, err := json.Marshal(order)
		if err != nil {
			t.Fatalf("encode order %q: %v", order.ID, err)
		}
		createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/orders", body)
		if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
			t.Fatalf("POST /examples/orders for %q status = %d, want 2xx; body: %s",
				order.ID, createResponse.StatusCode, createResponse.Body)
		}
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
		t.Fatalf("decode order list response: %v", err)
	}
	if len(listResult.Data) != len(expectedOrders) {
		t.Fatalf("order count = %d, want %d", len(listResult.Data), len(expectedOrders))
	}

	ordersByID := make(map[string]entity.OrderResponse, len(listResult.Data))
	for _, order := range listResult.Data {
		ordersByID[order.ID] = order
	}
	for _, expected := range expectedOrders {
		actual, ok := ordersByID[expected.ID]
		if !ok {
			t.Errorf("order %q not found in response", expected.ID)
			continue
		}
		if actual.UserID != expected.UserID ||
			actual.OrderNumber != expected.OrderNumber ||
			actual.Status != expected.Status ||
			actual.TotalAmount != expected.TotalAmount ||
			actual.IsActive != expected.IsActive {
			t.Errorf("order %q = %+v, want userID %q, orderNumber %q, status %q, totalAmount %.2f, isActive %t",
				expected.ID, actual, expected.UserID, expected.OrderNumber, expected.Status,
				expected.TotalAmount, expected.IsActive)
		}
		if actual.CreatedAt == "" {
			t.Errorf("order %q createdAt is empty", expected.ID)
		}
	}
}

func TestGetOrderByIDE2E(t *testing.T) {
	testServer, userID := newOrderTestServer(t)
	expected := entity.Order{
		ID:          "order-read-by-id-e2e-01",
		UserID:      userID,
		OrderNumber: "ORD-E2E-READ-ID-0001",
		Status:      "PROCESSING",
		TotalAmount: 98.75,
		IsActive:    true,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("encode order: %v", err)
	}

	createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/orders", body)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orders status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
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

	var result struct {
		Data entity.OrderResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &result); err != nil {
		t.Fatalf("decode order response: %v", err)
	}
	if result.Data.ID != expected.ID ||
		result.Data.UserID != expected.UserID ||
		result.Data.OrderNumber != expected.OrderNumber ||
		result.Data.Status != expected.Status ||
		result.Data.TotalAmount != expected.TotalAmount ||
		result.Data.IsActive != expected.IsActive {
		t.Fatalf("order response = %+v, want ID %q, userID %q, orderNumber %q, status %q, totalAmount %.2f, isActive %t",
			result.Data, expected.ID, expected.UserID, expected.OrderNumber, expected.Status,
			expected.TotalAmount, expected.IsActive)
	}
	if result.Data.CreatedAt == "" {
		t.Fatal("order createdAt is empty, want a formatted creation timestamp")
	}
}

func TestGetOrdersByUserIDE2E(t *testing.T) {
	testServer, userIDs := newOrderTestServerWithUsers(t, 2)
	expectedOrders := []entity.Order{
		{
			ID:          "order-user-read-e2e-01",
			UserID:      userIDs[0],
			OrderNumber: "ORD-E2E-USER-0001",
			Status:      "PENDING",
			TotalAmount: 35.50,
			IsActive:    true,
		},
		{
			ID:          "order-user-read-e2e-02",
			UserID:      userIDs[0],
			OrderNumber: "ORD-E2E-USER-0002",
			Status:      "FULFILLED",
			TotalAmount: 110.00,
			IsActive:    true,
		},
		{
			ID:          "order-user-read-e2e-03",
			UserID:      userIDs[1],
			OrderNumber: "ORD-E2E-OTHER-0001",
			Status:      "PENDING",
			TotalAmount: 999.00,
			IsActive:    true,
		},
	}
	for _, order := range expectedOrders {
		body, err := json.Marshal(order)
		if err != nil {
			t.Fatalf("encode order %q: %v", order.ID, err)
		}

		createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/orders", body)
		if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
			t.Fatalf("POST /examples/orders for %q status = %d, want 2xx; body: %s",
				order.ID, createResponse.StatusCode, createResponse.Body)
		}
	}

	listResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/users/"+userIDs[0]+"/orders",
		nil,
	)
	if listResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/users/%s/orders status = %d, want %d; body: %s",
			userIDs[0], listResponse.StatusCode, http.StatusOK, listResponse.Body)
	}

	var result struct {
		Data entity.OrderResponseArray `json:"data"`
	}
	if err := json.Unmarshal(listResponse.Body, &result); err != nil {
		t.Fatalf("decode user's order list response: %v", err)
	}
	if len(result.Data) != 2 {
		t.Fatalf("order count for user %q = %d, want 2", userIDs[0], len(result.Data))
	}

	ordersByID := make(map[string]entity.OrderResponse, len(result.Data))
	for _, order := range result.Data {
		if order.UserID != userIDs[0] {
			t.Errorf("order %q belongs to user %q, want only user %q orders",
				order.ID, order.UserID, userIDs[0])
		}
		ordersByID[order.ID] = order
	}
	for _, expected := range expectedOrders[:2] {
		actual, ok := ordersByID[expected.ID]
		if !ok {
			t.Errorf("order %q not found in user %q response", expected.ID, userIDs[0])
			continue
		}
		if actual.OrderNumber != expected.OrderNumber ||
			actual.Status != expected.Status ||
			actual.TotalAmount != expected.TotalAmount ||
			actual.IsActive != expected.IsActive {
			t.Errorf("order %q = %+v, want orderNumber %q, status %q, totalAmount %.2f, isActive %t",
				expected.ID, actual, expected.OrderNumber, expected.Status,
				expected.TotalAmount, expected.IsActive)
		}
	}
}

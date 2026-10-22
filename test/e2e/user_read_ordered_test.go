package e2e_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersOrderedByIDE2E(t *testing.T) {
	testServer, _ := newUserTestServer(t)
	path := "/examples/orm-queries/users/order-by"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty ordered-users response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty ordered-users data = %#v, want an empty array", emptyResult.Data)
	}

	users := []entity.User{
		{
			ID:           "user-order-e2e-b",
			Email:        "user-order-b@example.test",
			Password:     "plain-text-order-password-b",
			DepartmentID: "department-order-e2e",
			IsActive:     true,
		},
		{
			ID:           "user-order-e2e-a",
			Email:        "user-order-a@example.test",
			Password:     "plain-text-order-password-a",
			DepartmentID: "department-order-e2e",
			IsActive:     true,
		},
		{
			ID:           "user-order-e2e-c",
			Email:        "user-order-c@example.test",
			Password:     "plain-text-order-password-c",
			DepartmentID: "department-order-e2e",
			IsActive:     false,
		},
	}
	for _, user := range users {
		body, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("encode create user request: %v", err)
		}
		createResponse := sendRequest(
			t,
			testServer,
			http.MethodPost,
			"/examples/orm-queries/users",
			body,
		)
		if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
			t.Fatalf("POST /examples/orm-queries/users for %s status = %d, want 2xx; body: %s",
				user.ID, createResponse.StatusCode, createResponse.Body)
		}
	}

	response := sendRequest(t, testServer, http.MethodGet, path, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, response.StatusCode, http.StatusOK, response.Body)
	}
	var result struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode ordered-users response: %v", err)
	}

	expectedIDs := []string{"user-order-e2e-c", "user-order-e2e-b", "user-order-e2e-a"}
	actualIDs := make([]string, len(result.Data))
	for i, user := range result.Data {
		actualIDs[i] = user.ID
		if user.Email == "" || user.DepartmentID != "department-order-e2e" {
			t.Errorf("ordered user at index %d has unexpected fields: %+v", i, user)
		}
	}
	if !reflect.DeepEqual(actualIDs, expectedIDs) {
		t.Fatalf("ordered user IDs = %v, want descending order %v", actualIDs, expectedIDs)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode ordered-users response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode ordered-user fields: %v", err)
	}
	for i, user := range responseUsers {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("ordered user at index %d unexpectedly exposes password", i)
		}
	}
}

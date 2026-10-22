package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersE2E(t *testing.T) {
	testServer, _ := newUserTestServer(t)

	assertUserListResponse(t, sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orm-queries/users",
		nil,
	), map[string]entity.UserResponse{})

	expectedUsers := []entity.User{
		{
			ID:           "user-list-e2e-01",
			Email:        "user-list-01@example.test",
			Password:     "plain-text-list-password-01",
			DepartmentID: "department-list-e2e-01",
			IsActive:     true,
		},
		{
			ID:           "user-list-e2e-02",
			Email:        "user-list-02@example.test",
			Password:     "plain-text-list-password-02",
			DepartmentID: "department-list-e2e-02",
			IsActive:     false,
		},
	}
	for _, user := range expectedUsers {
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

	expectedByID := make(map[string]entity.UserResponse, len(expectedUsers))
	for _, user := range expectedUsers {
		expectedByID[user.ID] = entity.UserResponse{
			ID:           user.ID,
			Email:        user.Email,
			DepartmentID: user.DepartmentID,
			IsActive:     user.IsActive,
		}
	}
	assertUserListResponse(t, sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orm-queries/users",
		nil,
	), expectedByID)
}

func assertUserListResponse(t *testing.T, response testHTTPResponse, expectedByID map[string]entity.UserResponse) {
	t.Helper()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orm-queries/users status = %d, want %d; body: %s",
			response.StatusCode, http.StatusOK, response.Body)
	}

	var result struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode user list response: %v", err)
	}
	if result.Data == nil {
		t.Fatal("user list response data = null, want an array")
	}
	if len(result.Data) != len(expectedByID) {
		t.Fatalf("user list count = %d, want %d; body: %s",
			len(result.Data), len(expectedByID), response.Body)
	}

	for _, user := range result.Data {
		expected, found := expectedByID[user.ID]
		if !found {
			t.Fatalf("unexpected user in response: %+v", user)
		}
		if user.Email != expected.Email ||
			user.DepartmentID != expected.DepartmentID ||
			user.IsActive != expected.IsActive {
			t.Errorf("user response = %+v, want %+v", user, expected)
		}
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode response fields: %v", err)
	}
	var users []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &users); err != nil {
		t.Fatalf("decode user response fields: %v", err)
	}
	for i, user := range users {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("user at index %d unexpectedly exposes password", i)
		}
	}
}

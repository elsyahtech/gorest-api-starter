package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUserByIDE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	expected := entity.User{
		ID:           "user-by-id-e2e-01",
		Email:        "user-by-id@example.test",
		Password:     "plain-text-by-id-password",
		DepartmentID: "department-by-id-e2e-01",
		IsActive:     true,
	}
	query := app.Database().Table("users").Create(&expected)
	if query.Error != nil {
		t.Fatalf("seed user %s: %s: %v", expected.ID, query.Message, query.Error)
	}

	path := "/examples/orm-queries/users/" + expected.ID
	response := sendRequest(t, testServer, http.MethodGet, path, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, response.StatusCode, http.StatusOK, response.Body)
	}

	var result struct {
		Data entity.UserResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode user response: %v", err)
	}
	if result.Data.ID != expected.ID ||
		result.Data.Email != expected.Email ||
		result.Data.DepartmentID != expected.DepartmentID ||
		result.Data.IsActive != expected.IsActive {
		t.Fatalf("user response = %+v, want user with ID %q", result.Data, expected.ID)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode response fields: %v", err)
	}
	var responseUser map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUser); err != nil {
		t.Fatalf("decode user response fields: %v", err)
	}
	if _, hasPassword := responseUser["password"]; hasPassword {
		t.Error("user response unexpectedly exposes password")
	}

	missingID := "user-by-id-e2e-missing"
	missingPath := "/examples/orm-queries/users/" + missingID
	missingResponse := sendRequest(t, testServer, http.MethodGet, missingPath, nil)
	if missingResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			missingPath, missingResponse.StatusCode, http.StatusNotFound, missingResponse.Body)
	}
}

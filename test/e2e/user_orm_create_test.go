package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestCreateUserE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	expected := entity.User{
		ID:           "user-orm-create-e2e-01",
		Email:        "orm-create-e2e@example.test",
		Password:     "plain-text-orm-create-password",
		DepartmentID: "department-orm-create-e2e-01",
		IsActive:     true,
	}
	body, err := json.Marshal(expected)
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
		t.Fatalf("POST /examples/orm-queries/users status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode create user response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("create user response success = false, want true")
	}

	var persisted entity.User
	query := app.Database().Table("users").Where("id = ?", expected.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read persisted user: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != expected.ID ||
		persisted.Email != expected.Email ||
		persisted.DepartmentID != expected.DepartmentID ||
		persisted.IsActive != expected.IsActive {
		t.Fatalf("persisted user = %+v, does not match create request", persisted)
	}
	if persisted.Password == "" || persisted.Password == expected.Password {
		t.Fatal("persisted password is empty or plaintext; want a hashed password")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("persisted createdAt is zero, want a creation timestamp")
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orm-queries/users/"+expected.ID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orm-queries/users/%s status = %d, want %d; body: %s",
			expected.ID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}
	var getResult struct {
		Data entity.UserResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &getResult); err != nil {
		t.Fatalf("decode user response: %v", err)
	}
	if getResult.Data.ID != expected.ID ||
		getResult.Data.Email != expected.Email ||
		getResult.Data.DepartmentID != expected.DepartmentID ||
		getResult.Data.IsActive != expected.IsActive {
		t.Fatalf("user response = %+v, does not match create request", getResult.Data)
	}
	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(getResponse.Body, &responseFields); err != nil {
		t.Fatalf("decode user response fields: %v", err)
	}
	var dataFields map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &dataFields); err != nil {
		t.Fatalf("decode user response data fields: %v", err)
	}
	if _, hasPassword := dataFields["password"]; hasPassword {
		t.Fatal("user response unexpectedly exposes password")
	}

	duplicate := expected
	duplicate.ID = "user-orm-create-e2e-duplicate"
	duplicate.Password = "another-orm-create-password"
	duplicateBody, err := json.Marshal(duplicate)
	if err != nil {
		t.Fatalf("encode duplicate user request: %v", err)
	}
	duplicateResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users",
		duplicateBody,
	)
	if duplicateResponse.StatusCode < http.StatusBadRequest {
		t.Fatalf("POST duplicate user status = %d, want 4xx or 5xx; body: %s",
			duplicateResponse.StatusCode, duplicateResponse.Body)
	}

	var persistedAfterDuplicate entity.User
	query = app.Database().Table("users").Where("id = ?", expected.ID).Find(&persistedAfterDuplicate)
	if query.Error != nil {
		t.Fatalf("read original user after duplicate create: %s: %v", query.Message, query.Error)
	}
	if persistedAfterDuplicate.Password != persisted.Password {
		t.Fatal("failed duplicate create changed the original user's password")
	}

	var duplicateUser entity.User
	query = app.Database().Table("users").Where("id = ?", duplicate.ID).Find(&duplicateUser)
	if query.Error == nil {
		t.Fatalf("duplicate user was persisted unexpectedly: %+v", duplicateUser)
	}
}

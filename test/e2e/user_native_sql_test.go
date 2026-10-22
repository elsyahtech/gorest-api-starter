package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestCreateUserWithNativeSQLE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	expected := entity.User{
		ID:           "user-native-sql-e2e-01",
		Email:        "native-sql-e2e@example.test",
		Password:     "plain-text-native-sql-password",
		DepartmentID: "department-native-sql-e2e-01",
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
		"/examples/native-queries/users/sql",
		body,
	)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/native-queries/users/sql status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode native SQL create response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("native SQL create response success = false, want true")
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
		t.Fatalf("persisted user = %+v, does not match request", persisted)
	}
	if persisted.Password == "" || persisted.Password == expected.Password {
		t.Fatal("persisted password is empty or plaintext; want a hashed password")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("persisted createdAt is zero, want native SQL insert to store a timestamp")
	}

	duplicate := expected
	duplicate.ID = "user-native-sql-e2e-duplicate"
	duplicate.Password = "another-native-sql-password"
	duplicateBody, err := json.Marshal(duplicate)
	if err != nil {
		t.Fatalf("encode duplicate user request: %v", err)
	}
	duplicateResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/native-queries/users/sql",
		duplicateBody,
	)
	if duplicateResponse.StatusCode < http.StatusBadRequest {
		t.Fatalf("POST duplicate user status = %d, want 4xx or 5xx; body: %s",
			duplicateResponse.StatusCode, duplicateResponse.Body)
	}

	var afterFailure entity.User
	query = app.Database().Table("users").Where("id = ?", expected.ID).Find(&afterFailure)
	if query.Error != nil {
		t.Fatalf("read original user after failed insert: %s: %v", query.Message, query.Error)
	}
	if afterFailure.Password != persisted.Password {
		t.Fatal("failed native SQL insert changed the existing user's password")
	}

	var duplicateID entity.User
	query = app.Database().Table("users").Where("id = ?", duplicate.ID).Find(&duplicateID)
	if query.Error == nil {
		t.Fatalf("duplicate user was persisted unexpectedly: %+v", duplicateID)
	}
}

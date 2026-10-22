package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestCreateUsersInBulkE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	expectedUsers := entity.UserArray{
		{
			ID:           "user-bulk-e2e-01",
			Email:        "bulk-one@example.test",
			Password:     "bulk-password-one",
			DepartmentID: "department-bulk-01",
			IsActive:     true,
		},
		{
			ID:           "user-bulk-e2e-02",
			Email:        "bulk-two@example.test",
			Password:     "bulk-password-two",
			DepartmentID: "department-bulk-02",
			IsActive:     false,
		},
	}
	body, err := json.Marshal(expectedUsers)
	if err != nil {
		t.Fatalf("encode bulk user request: %v", err)
	}

	createResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users/bulk",
		body,
	)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orm-queries/users/bulk status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode bulk create response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("bulk create response success = false, want true")
	}

	for _, expected := range expectedUsers {
		var persisted entity.User
		query := app.Database().Table("users").Where("id = ?", expected.ID).Find(&persisted)
		if query.Error != nil {
			t.Fatalf("read persisted user %q: %s: %v", expected.ID, query.Message, query.Error)
		}
		if persisted.ID != expected.ID ||
			persisted.Email != expected.Email ||
			persisted.DepartmentID != expected.DepartmentID ||
			persisted.IsActive != expected.IsActive {
			t.Errorf("persisted user %q = %+v, does not match submitted fields", expected.ID, persisted)
		}
		if persisted.Password == "" || persisted.Password == expected.Password {
			t.Errorf("persisted password for user %q is empty or plaintext", expected.ID)
		}
		if persisted.CreatedAt.IsZero() {
			t.Errorf("persisted createdAt for user %q is zero", expected.ID)
		}
	}
}

func TestCreateUsersInBulkRejectsDuplicateRequestEmailE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	users := entity.UserArray{
		{
			ID:           "user-bulk-duplicate-request-01",
			Email:        "same-bulk-email@example.test",
			Password:     "bulk-password-one",
			DepartmentID: "department-bulk-01",
			IsActive:     true,
		},
		{
			ID:           "user-bulk-duplicate-request-02",
			Email:        "same-bulk-email@example.test",
			Password:     "bulk-password-two",
			DepartmentID: "department-bulk-02",
			IsActive:     true,
		},
	}
	body, err := json.Marshal(users)
	if err != nil {
		t.Fatalf("encode duplicate-email bulk request: %v", err)
	}

	response := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users/bulk",
		body,
	)
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("POST duplicate-email bulk request status = %d, want %d; body: %s",
			response.StatusCode, http.StatusBadRequest, response.Body)
	}
	assertUserCount(t, app, 0)
}

func TestCreateUsersInBulkRollsBackOnDatabaseErrorE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	existing := entity.User{
		ID:           "user-bulk-existing-01",
		Email:        "already-exists@example.test",
		Password:     "existing-password",
		DepartmentID: "department-bulk-existing",
		IsActive:     true,
	}
	query := app.Database().Table("users").Create(&existing)
	if query.Error != nil {
		t.Fatalf("seed existing user: %s: %v", query.Message, query.Error)
	}

	users := entity.UserArray{
		{
			ID:           "user-bulk-should-rollback-01",
			Email:        "new-before-conflict@example.test",
			Password:     "bulk-new-password",
			DepartmentID: "department-bulk-new",
			IsActive:     true,
		},
		{
			ID:           "user-bulk-conflict-01",
			Email:        existing.Email,
			Password:     "bulk-conflict-password",
			DepartmentID: "department-bulk-conflict",
			IsActive:     true,
		},
	}
	body, err := json.Marshal(users)
	if err != nil {
		t.Fatalf("encode database-conflict bulk request: %v", err)
	}

	response := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users/bulk",
		body,
	)
	if response.StatusCode < http.StatusBadRequest {
		t.Fatalf("POST database-conflict bulk request status = %d, want 4xx or 5xx; body: %s",
			response.StatusCode, response.Body)
	}

	var persistedExisting entity.User
	query = app.Database().Table("users").Where("id = ?", existing.ID).Find(&persistedExisting)
	if query.Error != nil {
		t.Fatalf("read existing user after failed batch: %s: %v", query.Message, query.Error)
	}
	if persistedExisting.ID != existing.ID || persistedExisting.Email != existing.Email {
		t.Fatalf("existing user after failed batch = %+v, want original record unchanged", persistedExisting)
	}

	assertUserCount(t, app, 1)
}

func assertUserCount(t *testing.T, app *gorest.App, want int) {
	t.Helper()

	var users entity.UserArray
	query := app.Database().Table("users").Find(&users)
	if query.Error != nil {
		t.Fatalf("count users: %s: %v", query.Message, query.Error)
	}
	if len(users) != want {
		t.Fatalf("user count = %d, want %d", len(users), want)
	}
}

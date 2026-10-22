package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestUpdateUserWithTransactionE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	initial := entity.User{
		ID:           "user-update-transaction-e2e-01",
		Email:        "update-transaction-old@example.test",
		Password:     "plain-text-update-transaction-password",
		DepartmentID: "department-update-transaction-old",
		IsActive:     true,
	}
	query := app.Database().Table("users").Create(&initial)
	if query.Error != nil {
		t.Fatalf("seed user %s: %s: %v", initial.ID, query.Message, query.Error)
	}

	updated := entity.User{
		ID:           initial.ID,
		Email:        "update-transaction-new@example.test",
		Password:     "plain-text-updated-transaction-password",
		DepartmentID: "department-update-transaction-new",
		IsActive:     false,
	}
	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode update user request: %v", err)
	}

	path := "/examples/orm-queries/users/transaction/" + initial.ID
	response := sendRequest(t, testServer, http.MethodPut, path, body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PUT %s status = %d, want 2xx; body: %s",
			path, response.StatusCode, response.Body)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode update user response: %v", err)
	}
	if !result.Success {
		t.Fatal("update user response success = false, want true")
	}

	var persisted entity.User
	query = app.Database().Table("users").Where("id = ?", initial.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read updated user: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != updated.ID ||
		persisted.Email != updated.Email ||
		persisted.DepartmentID != updated.DepartmentID ||
		persisted.IsActive != updated.IsActive {
		t.Fatalf("persisted user after update = %+v, want updated user %+v", persisted, updated)
	}
	if persisted.Password == "" || persisted.Password == updated.Password {
		t.Fatal("persisted password is empty or plaintext; want a hashed password")
	}
}

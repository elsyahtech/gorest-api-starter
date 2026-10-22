package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestDeleteUserWithTransactionE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	user := entity.User{
		ID:           "user-delete-transaction-e2e-01",
		Email:        "delete-transaction-e2e@example.test",
		Password:     "plain-text-delete-transaction-password",
		DepartmentID: "department-delete-transaction-e2e-01",
		IsActive:     true,
	}
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
		t.Fatalf("POST /examples/orm-queries/users status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var persisted entity.User
	query := app.Database().Table("users").Where("id = ?", user.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read user before delete: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != user.ID {
		t.Fatalf("persisted user ID = %q, want %q", persisted.ID, user.ID)
	}

	deleteResponse := sendRequest(
		t,
		testServer,
		http.MethodDelete,
		"/examples/orm-queries/users/transaction/"+user.ID,
		nil,
	)
	if deleteResponse.StatusCode < http.StatusOK || deleteResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("DELETE /examples/orm-queries/users/transaction/%s status = %d, want 2xx; body: %s",
			user.ID, deleteResponse.StatusCode, deleteResponse.Body)
	}

	var deleteResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(deleteResponse.Body, &deleteResult); err != nil {
		t.Fatalf("decode delete user response: %v", err)
	}
	if !deleteResult.Success {
		t.Fatal("delete user response success = false, want true")
	}

	var deletedUser entity.User
	query = app.Database().Table("users").Where("id = ?", user.ID).Find(&deletedUser)
	if query.Error == nil {
		t.Fatalf("user was persisted after transaction delete: %+v", deletedUser)
	}
}

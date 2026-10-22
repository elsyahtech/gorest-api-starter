package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestDeleteUserE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	userToDelete := entity.User{
		ID:           "user-delete-e2e-01",
		Email:        "delete-e2e@example.test",
		Password:     "plain-text-delete-password",
		DepartmentID: "department-delete-e2e-01",
		IsActive:     true,
	}
	userToKeep := entity.User{
		ID:           "user-delete-e2e-02",
		Email:        "keep-delete-e2e@example.test",
		Password:     "plain-text-keep-password",
		DepartmentID: "department-delete-e2e-02",
		IsActive:     true,
	}
	for _, user := range []entity.User{userToDelete, userToKeep} {
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

	deleteResponse := sendRequest(
		t,
		testServer,
		http.MethodDelete,
		"/examples/orm-queries/users/"+userToDelete.ID,
		nil,
	)
	if deleteResponse.StatusCode < http.StatusOK || deleteResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("DELETE /examples/orm-queries/users/%s status = %d, want 2xx; body: %s",
			userToDelete.ID, deleteResponse.StatusCode, deleteResponse.Body)
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
	query := app.Database().Table("users").Where("id = ?", userToDelete.ID).Find(&deletedUser)
	if query.Error == nil {
		t.Fatalf("user remains after ORM delete: %+v", deletedUser)
	}

	var remainingUser entity.User
	query = app.Database().Table("users").Where("id = ?", userToKeep.ID).Find(&remainingUser)
	if query.Error != nil {
		t.Fatalf("read unaffected user after delete: %s: %v", query.Message, query.Error)
	}
	if remainingUser.ID != userToKeep.ID || remainingUser.Email != userToKeep.Email {
		t.Fatalf("unaffected user after delete = %+v, want ID %q and email %q",
			remainingUser, userToKeep.ID, userToKeep.Email)
	}
}

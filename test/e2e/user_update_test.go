package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

func TestUpdateUserE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	createdAt := time.Date(2025, time.January, 15, 10, 0, 0, 0, time.UTC)
	initial := entity.User{
		CreatedAt:    createdAt,
		ID:           "user-update-e2e-01",
		Email:        "update-old@example.test",
		Password:     "original-hashed-password",
		DepartmentID: "department-update-old",
		IsActive:     true,
	}
	ensureE2EUserTimestamp(&initial)
	query := app.Database().Table("users").Create(&initial)
	if query.Error != nil {
		t.Fatalf("seed user %s: %s: %v", initial.ID, query.Message, query.Error)
	}

	updated := entity.User{
		ID:           initial.ID,
		Email:        "update-new@example.test",
		Password:     "plain-text-update-password",
		DepartmentID: "department-update-new",
		IsActive:     false,
	}
	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode update user request: %v", err)
	}

	path := "/examples/orm-queries/users/" + initial.ID
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
		t.Fatalf("persisted user after update = %+v", persisted)
	}
	if persisted.Password == "" || persisted.Password == updated.Password {
		t.Fatal("persisted password is empty or plaintext; want a hashed password")
	}
	matchesPassword, message, _, err := helper.VerifyArgon(updated.Password, persisted.Password)
	if err != nil {
		t.Fatalf("verify persisted password hash: %s: %v", message, err)
	}
	if !matchesPassword {
		t.Fatal("persisted password hash does not match the submitted password")
	}
	if !persisted.CreatedAt.Equal(createdAt) {
		t.Fatalf("persisted createdAt = %s, want unchanged value %s",
			persisted.CreatedAt.Format(time.RFC3339Nano), createdAt.Format(time.RFC3339Nano))
	}

	missingID := "user-update-e2e-missing"
	missingBody, err := json.Marshal(entity.User{
		Email:        "update-missing@example.test",
		Password:     "plain-text-missing-user-password",
		DepartmentID: "department-update-missing",
	})
	if err != nil {
		t.Fatalf("encode missing user update request: %v", err)
	}
	missingPath := "/examples/orm-queries/users/" + missingID
	missingResponse := sendRequest(t, testServer, http.MethodPut, missingPath, missingBody)
	if missingResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("PUT %s status = %d, want %d; body: %s",
			missingPath, missingResponse.StatusCode, http.StatusNotFound, missingResponse.Body)
	}
}

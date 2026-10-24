package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestUpdateUserWithSelectedColumnsE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	initialCreatedAt := time.Date(2024, time.December, 5, 9, 15, 0, 0, time.UTC)
	initialUpdatedAt := time.Date(2025, time.January, 10, 12, 0, 0, 0, time.UTC)
	initial := entity.User{
		CreatedAt:    initialCreatedAt,
		ID:           "user-update-selected-e2e-01",
		Email:        "update-selected-old@example.test",
		Password:     "original-hashed-password",
		DepartmentID: "department-update-selected-original",
		IsActive:     true,
		UpdatedAt:    &initialUpdatedAt,
	}
	ensureE2EUserTimestamp(&initial)
	query := app.Database().Table("users").Create(&initial)
	if query.Error != nil {
		t.Fatalf("seed user %s: %s: %v", initial.ID, query.Message, query.Error)
	}

	expectedUpdatedAt := time.Date(2025, time.February, 20, 15, 30, 0, 0, time.UTC)
	updated := entity.User{
		ID:           initial.ID,
		Email:        "update-selected-new@example.test",
		Password:     "should-not-update-password",
		DepartmentID: "department-update-selected-should-not-update",
		IsActive:     false,
		UpdatedAt:    &expectedUpdatedAt,
	}
	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode update user request: %v", err)
	}

	path := "/examples/orm-queries/users/selected-columns/" + initial.ID
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
	if persisted.ID != initial.ID || persisted.Email != updated.Email {
		t.Fatalf("persisted user identity and email = (%q, %q), want (%q, %q)",
			persisted.ID, persisted.Email, initial.ID, updated.Email)
	}
	if persisted.Password != initial.Password ||
		persisted.DepartmentID != initial.DepartmentID ||
		persisted.IsActive != initial.IsActive {
		t.Fatalf("unselected fields changed: password=%q departmentID=%q isActive=%t",
			persisted.Password, persisted.DepartmentID, persisted.IsActive)
	}
	if persisted.UpdatedAt == nil || !persisted.UpdatedAt.Equal(expectedUpdatedAt) {
		t.Fatalf("persisted updatedAt = %v, want %s",
			persisted.UpdatedAt, expectedUpdatedAt.Format(time.RFC3339Nano))
	}
	if !persisted.CreatedAt.Equal(initialCreatedAt) {
		t.Fatalf("persisted createdAt = %s, want unchanged value %s",
			persisted.CreatedAt.Format(time.RFC3339Nano), initialCreatedAt.Format(time.RFC3339Nano))
	}

	missingID := "user-update-selected-e2e-missing"
	missingBody, err := json.Marshal(entity.User{
		ID:    missingID,
		Email: "update-selected-missing@example.test",
	})
	if err != nil {
		t.Fatalf("encode missing user update request: %v", err)
	}
	missingPath := "/examples/orm-queries/users/selected-columns/" + missingID
	missingResponse := sendRequest(t, testServer, http.MethodPut, missingPath, missingBody)
	if missingResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("PUT %s status = %d, want %d; body: %s",
			missingPath, missingResponse.StatusCode, http.StatusNotFound, missingResponse.Body)
	}
}

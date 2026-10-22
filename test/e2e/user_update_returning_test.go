package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestUpdateUserReturningResultE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	initialCreatedAt := time.Date(2024, time.December, 5, 9, 15, 0, 0, time.UTC)
	initial := entity.User{
		CreatedAt:    initialCreatedAt,
		ID:           "user-update-returning-e2e-01",
		Email:        "update-returning-old@example.test",
		Password:     "original-hashed-password",
		DepartmentID: "department-update-returning-old",
		IsActive:     true,
	}
	query := app.Database().Table("users").Create(&initial)
	if query.Error != nil {
		t.Fatalf("seed user %s: %s: %v", initial.ID, query.Message, query.Error)
	}

	updated := entity.User{
		ID:           initial.ID,
		Email:        "update-returning-new@example.test",
		Password:     "plain-text-update-returning-password",
		DepartmentID: "department-update-returning-new",
		IsActive:     false,
	}
	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode update user request: %v", err)
	}

	path := "/examples/orm-queries/users/returning/" + initial.ID
	response := sendRequest(t, testServer, http.MethodPut, path, body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PUT %s status = %d, want 2xx; body: %s",
			path, response.StatusCode, response.Body)
	}

	var result struct {
		Success bool                `json:"success"`
		Data    entity.UserResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode update user response: %v", err)
	}
	if !result.Success {
		t.Fatal("update user response success = false, want true")
	}
	if result.Data.ID != updated.ID || result.Data.Email != updated.Email {
		t.Fatalf("returned user = %+v, want ID %q and email %q",
			result.Data, updated.ID, updated.Email)
	}
	if result.Data.CreatedAt == "" {
		t.Fatal("returned createdAt is empty")
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode response fields: %v", err)
	}
	var responseUser map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUser); err != nil {
		t.Fatalf("decode returned user fields: %v", err)
	}
	for _, field := range []string{"id", "email", "createdAt"} {
		if _, found := responseUser[field]; !found {
			t.Errorf("update response does not include returned field %q", field)
		}
	}
	if _, hasPassword := responseUser["password"]; hasPassword {
		t.Error("update response unexpectedly exposes password")
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
	if !persisted.CreatedAt.Equal(initialCreatedAt) {
		t.Fatalf("persisted createdAt = %s, want unchanged value %s",
			persisted.CreatedAt.Format(time.RFC3339Nano), initialCreatedAt.Format(time.RFC3339Nano))
	}

	missingID := "user-update-returning-e2e-missing"
	missingBody, err := json.Marshal(entity.User{
		ID:       missingID,
		Email:    "update-returning-missing@example.test",
		Password: "plain-text-missing-user-password",
	})
	if err != nil {
		t.Fatalf("encode missing user update request: %v", err)
	}
	missingPath := "/examples/orm-queries/users/returning/" + missingID
	missingResponse := sendRequest(t, testServer, http.MethodPut, missingPath, missingBody)
	if missingResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("PUT %s status = %d, want %d; body: %s",
			missingPath, missingResponse.StatusCode, http.StatusNotFound, missingResponse.Body)
	}
}

package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

func TestUpdateUserWithNativeSQLE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	createdAt := time.Date(2025, time.January, 15, 10, 0, 0, 0, time.UTC)
	initial := entity.User{
		CreatedAt:    createdAt,
		ID:           "user-update-native-sql-e2e-01",
		Email:        "update-native-sql-old@example.test",
		Password:     "original-hashed-password",
		DepartmentID: "department-update-native-sql-old",
		IsActive:     true,
	}
	ensureE2EUserTimestamp(&initial)
	query := app.Database().Table("users").Create(&initial)
	if query.Error != nil {
		t.Fatalf("seed user %s: %s: %v", initial.ID, query.Message, query.Error)
	}

	updated := entity.User{
		Email:        "update-native-sql-new@example.test",
		Password:     "plain-text-native-sql-updated-password",
		DepartmentID: "department-update-native-sql-new",
	}
	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode native SQL update request: %v", err)
	}

	path := "/examples/native-queries/users/sql/" + initial.ID
	response := sendRequest(t, testServer, http.MethodPut, path, body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PUT %s status = %d, want 2xx; body: %s",
			path, response.StatusCode, response.Body)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode native SQL update response: %v", err)
	}
	if !result.Success {
		t.Fatal("native SQL update response success = false, want true")
	}

	var persisted entity.User
	query = app.Database().Table("users").Where("id = ?", initial.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read updated user: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != initial.ID ||
		persisted.Email != updated.Email ||
		persisted.DepartmentID != updated.DepartmentID ||
		!persisted.IsActive {
		t.Fatalf("persisted user after native SQL update = %+v", persisted)
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
	if persisted.UpdatedAt == nil || persisted.UpdatedAt.IsZero() {
		t.Fatal("persisted updatedAt is nil or zero")
	}
}

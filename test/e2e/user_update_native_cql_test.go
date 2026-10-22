package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

func TestUpdateUserWithNativeCQLE2E(t *testing.T) {
	testServer, app := newScyllaUserTestServer(t)
	createdAt := time.Date(2025, time.January, 15, 10, 0, 0, 0, time.UTC)
	initial := entity.User{
		ID:           "user-update-native-cql-e2e-01",
		Email:        "update-native-cql-old@example.test",
		Password:     "original-hashed-password",
		DepartmentID: "department-update-native-cql-old",
		IsActive:     true,
		CreatedAt:    createdAt,
	}
	if err := executeScyllaCQL(t, app, `
		INSERT INTO users (id, email, password, department_id, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, initial.ID, initial.Email, initial.Password, initial.DepartmentID, initial.IsActive, initial.CreatedAt); err != nil {
		t.Fatalf("seed ScyllaDB user: %v", err)
	}

	updated := entity.User{
		ID:           initial.ID,
		Email:        "update-native-cql-new@example.test",
		Password:     "plain-text-native-cql-updated-password",
		DepartmentID: "department-update-native-cql-new",
		IsActive:     false,
	}
	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode native CQL update request: %v", err)
	}
	path := "/examples/native-queries/users/scylla/" + initial.ID
	response := sendRequest(t, testServer, http.MethodPut, path, body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PUT %s status = %d, want 2xx; body: %s",
			path, response.StatusCode, response.Body)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode native CQL update response: %v", err)
	}
	if !result.Success {
		t.Fatal("native CQL update response success = false, want true")
	}

	ctx, cancel := app.NewContext(10 * time.Second)
	defer cancel()
	iter, _, _, err := app.Database().QueryCQL(ctx, `
		SELECT id, email, password, department_id, is_active, created_at, updated_at
		FROM users WHERE id = ?
	`, initial.ID)
	if err != nil {
		t.Fatalf("query updated ScyllaDB user: %v", err)
	}
	var persisted entity.User
	if !iter.Scan(
		&persisted.ID,
		&persisted.Email,
		&persisted.Password,
		&persisted.DepartmentID,
		&persisted.IsActive,
		&persisted.CreatedAt,
		&persisted.UpdatedAt,
	) {
		if err := iter.Close(); err != nil {
			t.Fatalf("close ScyllaDB user query: %v", err)
		}
		t.Fatal("ScyllaDB user query returned no row")
	}
	if err := iter.Close(); err != nil {
		t.Fatalf("close ScyllaDB user query: %v", err)
	}
	if persisted.ID != updated.ID ||
		persisted.Email != updated.Email ||
		persisted.DepartmentID != updated.DepartmentID ||
		persisted.IsActive != updated.IsActive {
		t.Fatalf("persisted ScyllaDB user after update = %+v", persisted)
	}
	if persisted.Password == "" || persisted.Password == updated.Password {
		t.Fatal("persisted ScyllaDB password is empty or plaintext; want a hashed password")
	}
	matchesPassword, message, _, err := helper.VerifyArgon(updated.Password, persisted.Password)
	if err != nil {
		t.Fatalf("verify persisted ScyllaDB password hash: %s: %v", message, err)
	}
	if !matchesPassword {
		t.Fatal("persisted ScyllaDB password hash does not match submitted password")
	}
	if !persisted.CreatedAt.Equal(createdAt) {
		t.Fatalf("persisted ScyllaDB createdAt = %s, want unchanged value %s",
			persisted.CreatedAt.Format(time.RFC3339Nano), createdAt.Format(time.RFC3339Nano))
	}
	if persisted.UpdatedAt == nil || persisted.UpdatedAt.IsZero() {
		t.Fatal("persisted ScyllaDB updatedAt is nil or zero")
	}

	missingID := "user-update-native-cql-e2e-missing"
	missingBody, err := json.Marshal(entity.User{
		ID:       missingID,
		Email:    "update-native-cql-missing@example.test",
		Password: "plain-text-missing-user-password",
	})
	if err != nil {
		t.Fatalf("encode missing ScyllaDB user update request: %v", err)
	}
	missingPath := "/examples/native-queries/users/scylla/" + missingID
	missingResponse := sendRequest(t, testServer, http.MethodPut, missingPath, missingBody)
	if missingResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("PUT %s status = %d, want %d; body: %s",
			missingPath, missingResponse.StatusCode, http.StatusNotFound, missingResponse.Body)
	}
}

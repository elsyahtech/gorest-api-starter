package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestPatchUserActiveStatusE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	createdAt := time.Date(2025, time.January, 15, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2025, time.January, 16, 10, 0, 0, 0, time.UTC)
	initial := entity.User{
		CreatedAt:    createdAt,
		UpdatedAt:    &updatedAt,
		ID:           "user-patch-active-e2e-01",
		Email:        "patch-active@example.test",
		Password:     "original-hashed-password",
		DepartmentID: "department-patch-active",
		IsActive:     true,
	}
	ensureE2EUserTimestamp(&initial)
	query := app.Database().Table("users").Create(&initial)
	if query.Error != nil {
		t.Fatalf("seed user %s: %s: %v", initial.ID, query.Message, query.Error)
	}

	path := "/examples/orm-queries/users/" + initial.ID
	response := sendRequest(t, testServer, http.MethodPatch, path, []byte(`{"isActive":false}`))
	if response.StatusCode != http.StatusOK {
		t.Fatalf("PATCH %s status = %d, want %d; body: %s",
			path, response.StatusCode, http.StatusOK, response.Body)
	}
	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode patch user response: %v", err)
	}
	if !result.Success {
		t.Fatal("patch user response success = false, want true")
	}

	var persisted entity.User
	query = app.Database().Table("users").Where("id = ?", initial.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read patched user: %s: %v", query.Message, query.Error)
	}
	if persisted.IsActive {
		t.Fatal("persisted isActive = true, want false")
	}
	if persisted.ID != initial.ID ||
		persisted.Email != initial.Email ||
		persisted.Password != initial.Password ||
		persisted.DepartmentID != initial.DepartmentID {
		t.Fatalf("PATCH changed fields outside isActive: %+v", persisted)
	}
	if !persisted.CreatedAt.Equal(createdAt) {
		t.Fatalf("persisted createdAt = %s, want unchanged value %s",
			persisted.CreatedAt.Format(time.RFC3339Nano), createdAt.Format(time.RFC3339Nano))
	}
	if persisted.UpdatedAt == nil || !persisted.UpdatedAt.After(updatedAt) {
		t.Fatalf("persisted updatedAt = %v, want a timestamp after %s",
			persisted.UpdatedAt, updatedAt.Format(time.RFC3339Nano))
	}

	emptyResponse := sendRequest(t, testServer, http.MethodPatch, path, []byte(`{}`))
	if emptyResponse.StatusCode != http.StatusBadRequest {
		t.Fatalf("PATCH %s without isActive status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusBadRequest, emptyResponse.Body)
	}

	missingPath := "/examples/orm-queries/users/user-patch-active-e2e-missing"
	missingResponse := sendRequest(
		t,
		testServer,
		http.MethodPatch,
		missingPath,
		[]byte(`{"isActive":false}`),
	)
	if missingResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("PATCH %s status = %d, want %d; body: %s",
			missingPath, missingResponse.StatusCode, http.StatusNotFound, missingResponse.Body)
	}
}

package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestUpdateUsersInBulkE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	createdAt := time.Date(2025, time.January, 15, 10, 0, 0, 0, time.UTC)
	initialUsers := entity.UserArray{
		{
			CreatedAt:    createdAt,
			ID:           "user-bulk-update-e2e-01",
			Email:        "bulk-update-old-01@example.test",
			Password:     "original-password-01",
			DepartmentID: "department-bulk-update-old-01",
			IsActive:     true,
		},
		{
			CreatedAt:    createdAt,
			ID:           "user-bulk-update-e2e-02",
			Email:        "bulk-update-old-02@example.test",
			Password:     "original-password-02",
			DepartmentID: "department-bulk-update-old-02",
			IsActive:     true,
		},
	}
	for i := range initialUsers {
		ensureE2EUserTimestamp(&initialUsers[i])
		query := app.Database().Table("users").Create(&initialUsers[i])
		if query.Error != nil {
			t.Fatalf("seed user %s: %s: %v", initialUsers[i].ID, query.Message, query.Error)
		}
	}

	updatedUsers := entity.UserArray{
		{
			ID:           "user-bulk-update-e2e-01",
			Email:        "bulk-update-new-01@example.test",
			Password:     "updated-password-01",
			DepartmentID: "department-bulk-update-new-01",
			IsActive:     false,
		},
		{
			ID:           "user-bulk-update-e2e-02",
			Email:        "bulk-update-new-02@example.test",
			Password:     "updated-password-02",
			DepartmentID: "department-bulk-update-new-02",
			IsActive:     false,
		},
	}
	body, err := json.Marshal(updatedUsers)
	if err != nil {
		t.Fatalf("encode bulk update request: %v", err)
	}

	path := "/examples/orm-queries/users/bulk"
	response := sendRequest(t, testServer, http.MethodPut, path, body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PUT %s status = %d, want 2xx; body: %s",
			path, response.StatusCode, response.Body)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode bulk update response: %v", err)
	}
	if !result.Success {
		t.Fatal("bulk update response success = false, want true")
	}

	for i, expected := range updatedUsers {
		var persisted entity.User
		query := app.Database().Table("users").Where("id = ?", expected.ID).Find(&persisted)
		if query.Error != nil {
			t.Fatalf("read updated user %q: %s: %v", expected.ID, query.Message, query.Error)
		}
		if persisted.ID != expected.ID ||
			persisted.Email != expected.Email ||
			persisted.DepartmentID != expected.DepartmentID ||
			persisted.IsActive != expected.IsActive {
			t.Errorf("persisted user %q = %+v, does not match update request", expected.ID, persisted)
		}
		if persisted.Password == "" || persisted.Password == expected.Password {
			t.Errorf("persisted password for user %q is empty or plaintext", expected.ID)
		}
		if persisted.UpdatedAt == nil || persisted.UpdatedAt.IsZero() {
			t.Errorf("persisted updatedAt for user %q is nil or zero", expected.ID)
		}
		if !persisted.CreatedAt.Equal(createdAt) {
			t.Errorf("persisted createdAt for user %q = %s, want unchanged value %s",
				expected.ID, persisted.CreatedAt.Format(time.RFC3339Nano), createdAt.Format(time.RFC3339Nano))
		}
		if persisted.Password == initialUsers[i].Password {
			t.Errorf("persisted password for user %q was not updated", expected.ID)
		}
	}
}

func TestUpdateUsersInBulkRejectsMissingUserE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	existing := entity.User{
		ID:           "user-bulk-update-existing",
		Email:        "bulk-update-existing@example.test",
		Password:     "existing-password",
		DepartmentID: "department-bulk-update-existing",
		IsActive:     true,
	}
	ensureE2EUserTimestamp(&existing)
	query := app.Database().Table("users").Create(&existing)
	if query.Error != nil {
		t.Fatalf("seed existing user: %s: %v", query.Message, query.Error)
	}

	users := entity.UserArray{
		{
			ID:           existing.ID,
			Email:        "bulk-update-existing-changed@example.test",
			Password:     "changed-existing-password",
			DepartmentID: "dept-bulk-update-existing-changed",
			IsActive:     false,
		},
		{
			ID:           "user-bulk-update-missing",
			Email:        "bulk-update-missing@example.test",
			Password:     "missing-user-password",
			DepartmentID: "department-bulk-update-missing",
			IsActive:     false,
		},
	}
	body, err := json.Marshal(users)
	if err != nil {
		t.Fatalf("encode missing-user bulk update request: %v", err)
	}

	path := "/examples/orm-queries/users/bulk"
	response := sendRequest(t, testServer, http.MethodPut, path, body)
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("PUT %s with a missing user status = %d, want %d; body: %s",
			path, response.StatusCode, http.StatusNotFound, response.Body)
	}

	var persisted entity.User
	query = app.Database().Table("users").Where("id = ?", existing.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read existing user after rejected bulk update: %s: %v", query.Message, query.Error)
	}
	if persisted.Email != existing.Email ||
		persisted.Password != existing.Password ||
		persisted.DepartmentID != existing.DepartmentID ||
		persisted.IsActive != existing.IsActive {
		t.Fatalf("existing user changed despite missing ID rejection: %+v", persisted)
	}
}

func TestUpdateUsersInBulkRejectsDuplicatePrimaryKeyE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	existing := entity.User{
		ID:           "user-bulk-update-duplicate",
		Email:        "bulk-update-duplicate-original@example.test",
		Password:     "original-password",
		DepartmentID: "dept-bulk-dup-original",
		IsActive:     true,
	}
	ensureE2EUserTimestamp(&existing)
	query := app.Database().Table("users").Create(&existing)
	if query.Error != nil {
		t.Fatalf("seed user: %s: %v", query.Message, query.Error)
	}

	users := entity.UserArray{
		{
			ID:           existing.ID,
			Email:        "bulk-update-duplicate-first@example.test",
			Password:     "first-updated-password",
			DepartmentID: "dept-bulk-dup-first",
			IsActive:     false,
		},
		{
			ID:           existing.ID,
			Email:        "bulk-update-duplicate-second@example.test",
			Password:     "second-updated-password",
			DepartmentID: "dept-bulk-dup-second",
			IsActive:     false,
		},
	}
	body, err := json.Marshal(users)
	if err != nil {
		t.Fatalf("encode duplicate-primary-key bulk update request: %v", err)
	}

	path := "/examples/orm-queries/users/bulk"
	response := sendRequest(t, testServer, http.MethodPut, path, body)
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("PUT %s with duplicate primary keys status = %d, want %d; body: %s",
			path, response.StatusCode, http.StatusBadRequest, response.Body)
	}

	var persisted entity.User
	query = app.Database().Table("users").Where("id = ?", existing.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read user after rejected duplicate update: %s: %v", query.Message, query.Error)
	}
	if persisted.Email != existing.Email ||
		persisted.Password != existing.Password ||
		persisted.DepartmentID != existing.DepartmentID ||
		persisted.IsActive != existing.IsActive {
		t.Fatalf("user changed despite duplicate primary-key rejection: %+v", persisted)
	}
}

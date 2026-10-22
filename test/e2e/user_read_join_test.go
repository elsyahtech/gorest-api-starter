package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersWithDepartmentE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, message, _, err := app.Database().ExecSQL(ctx, `
		CREATE TABLE departments (
			id VARCHAR(36) NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT NULL,
			deleted_at DATETIME DEFAULT NULL
		)
	`)
	if err != nil {
		t.Fatalf("create departments table for join E2E: %s: %v", message, err)
	}
	path := "/examples/orm-queries/users/join"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.JoinResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty joined-users response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty joined-users data = %#v, want an empty array", emptyResult.Data)
	}

	departments := []entity.Department{
		{
			ID:          "department-join-e2e-01",
			Name:        "Engineering",
			Description: "Builds and maintains software",
			IsActive:    true,
		},
	}
	for _, department := range departments {
		query := app.Database().Table("departments").Create(&department)
		if query.Error != nil {
			t.Fatalf("seed department %s: %s: %v", department.ID, query.Message, query.Error)
		}
	}

	users := []entity.User{
		{
			ID:           "user-join-e2e-01",
			Email:        "joined-user@example.test",
			Password:     "test-password-joined",
			DepartmentID: departments[0].ID,
			IsActive:     true,
		},
		{
			ID:           "user-join-e2e-02",
			Email:        "orphan-user@example.test",
			Password:     "test-password-orphan",
			DepartmentID: "department-join-e2e-missing",
			IsActive:     false,
		},
	}
	for _, user := range users {
		query := app.Database().Table("users").Create(&user)
		if query.Error != nil {
			t.Fatalf("seed user %s: %s: %v", user.ID, query.Message, query.Error)
		}
	}

	response := sendRequest(t, testServer, http.MethodGet, path, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, response.StatusCode, http.StatusOK, response.Body)
	}
	var result struct {
		Data entity.JoinResponseArray `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode joined-users response: %v", err)
	}
	if len(result.Data) != len(users) {
		t.Fatalf("joined user count = %d, want %d; data=%+v", len(result.Data), len(users), result.Data)
	}

	usersByID := make(map[string]entity.JoinResponse, len(result.Data))
	for _, user := range result.Data {
		usersByID[user.ID] = user
	}
	joined := usersByID[users[0].ID]
	if joined.Email != users[0].Email ||
		joined.DepartmentID != departments[0].ID ||
		joined.DepartmentName != departments[0].Name ||
		joined.DepartmentDescription != departments[0].Description ||
		!joined.IsActive {
		t.Fatalf("joined department data = %+v, want user and department fields from fixture", joined)
	}

	orphan := usersByID[users[1].ID]
	if orphan.Email != users[1].Email ||
		orphan.DepartmentID != users[1].DepartmentID ||
		orphan.DepartmentName != "" ||
		orphan.DepartmentDescription != "" ||
		orphan.IsActive {
		t.Fatalf("left join result for user without department = %+v", orphan)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode joined-users response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode joined-user fields: %v", err)
	}
	for i, user := range responseUsers {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("joined user at index %d unexpectedly exposes password", i)
		}
	}
}

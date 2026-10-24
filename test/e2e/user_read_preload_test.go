package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersWithPreloadedDepartmentE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	path := "/examples/orm-queries/users/preload"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty preloaded-users response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty preloaded-users data = %#v, want an empty array", emptyResult.Data)
	}

	department := entity.Department{
		ID:          "department-preload-e2e-01",
		Name:        "Research",
		Description: "Researches and prototypes new ideas",
		IsActive:    true,
	}
	ensureE2EDepartmentTimestamp(&department)
	departmentQuery := app.Database().Table("departments").Create(&department)
	if departmentQuery.Error != nil {
		t.Fatalf("seed department: %s: %v", departmentQuery.Message, departmentQuery.Error)
	}

	users := []entity.User{
		{
			ID:           "user-preload-e2e-01",
			Email:        "preloaded-user@example.test",
			Password:     "test-password-preloaded",
			DepartmentID: department.ID,
			IsActive:     true,
		},
		{
			ID:           "user-preload-e2e-02",
			Email:        "unmatched-user@example.test",
			Password:     "test-password-unmatched",
			DepartmentID: "department-preload-e2e-missing",
			IsActive:     false,
		},
	}
	for _, user := range users {
		ensureE2EUserTimestamp(&user)
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
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode preloaded-users response: %v", err)
	}
	if len(result.Data) != len(users) {
		t.Fatalf("preloaded user count = %d, want %d; data=%+v",
			len(result.Data), len(users), result.Data)
	}

	usersByID := make(map[string]entity.UserResponse, len(result.Data))
	for _, user := range result.Data {
		usersByID[user.ID] = user
	}
	preloaded := usersByID[users[0].ID]
	if preloaded.Department == nil {
		t.Fatalf("user %q has no preloaded department", preloaded.ID)
	}
	if preloaded.Department.ID != department.ID ||
		preloaded.Department.Name != department.Name ||
		preloaded.Department.Description != department.Description ||
		!preloaded.Department.IsActive {
		t.Fatalf("preloaded department = %+v, want %+v", preloaded.Department, department)
	}

	unmatched := usersByID[users[1].ID]
	if unmatched.Department != nil {
		t.Fatalf("user without matching department has relation %+v, want nil", unmatched.Department)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode preloaded response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode preloaded user fields: %v", err)
	}
	for i, user := range responseUsers {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("preloaded user at index %d unexpectedly exposes password", i)
		}
	}
}

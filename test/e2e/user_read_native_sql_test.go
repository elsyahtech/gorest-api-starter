package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersWithNativeSQLE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	path := "/examples/native-queries/users/sql"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty native SQL users response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty native SQL users data = %#v, want an empty array", emptyResult.Data)
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	users := []entity.User{
		{
			ID:           "user-native-sql-read-e2e-01",
			Email:        "native-sql-read-01@example.test",
			Password:     "hashed-test-password-01",
			DepartmentID: "department-native-sql-read-e2e",
			IsActive:     true,
			CreatedAt:    now,
		},
		{
			ID:           "user-native-sql-read-e2e-02",
			Email:        "native-sql-read-02@example.test",
			Password:     "hashed-test-password-02",
			DepartmentID: "department-native-sql-read-e2e",
			IsActive:     false,
			CreatedAt:    now,
		},
		{
			ID:           "user-native-sql-read-e2e-03",
			Email:        "native-sql-read-03@example.test",
			Password:     "hashed-test-password-03",
			DepartmentID: "department-native-sql-read-e2e",
			IsActive:     true,
			CreatedAt:    now,
		},
	}
	for _, user := range users {
		query := app.Database().Table("users").Create(&user)
		if query.Error != nil {
			t.Fatalf("seed user %s: %s: %v", user.ID, query.Message, query.Error)
		}
	}

	deletedAt := now.Add(time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	deletedAtPlaceholder := e2eSQLPlaceholder(1)
	userIDPlaceholder := e2eSQLPlaceholder(2)
	_, message, _, err := app.Database().ExecSQL(ctx,
		"UPDATE users SET deleted_at = "+deletedAtPlaceholder+" WHERE id = "+userIDPlaceholder,
		deletedAt,
		users[2].ID,
	)
	if err != nil {
		t.Fatalf("soft-delete native SQL fixture: %s: %v", message, err)
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
		t.Fatalf("decode native SQL users response: %v", err)
	}
	if len(result.Data) != 2 {
		t.Fatalf("native SQL user count = %d, want 2; data=%+v", len(result.Data), result.Data)
	}

	usersByID := make(map[string]entity.UserResponse, len(result.Data))
	for _, user := range result.Data {
		usersByID[user.ID] = user
	}
	for _, expected := range users[:2] {
		user, found := usersByID[expected.ID]
		if !found {
			t.Errorf("native SQL result is missing user %q", expected.ID)
			continue
		}
		if user.Email != expected.Email ||
			user.DepartmentID != expected.DepartmentID ||
			user.IsActive != expected.IsActive ||
			user.CreatedAt == "" ||
			user.DeletedAt != nil {
			t.Errorf("native SQL user response = %+v, does not match active fixture %+v",
				user, expected)
		}
	}
	if _, found := usersByID[users[2].ID]; found {
		t.Fatalf("soft-deleted user %q was returned by native SQL query", users[2].ID)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode native SQL response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode native SQL user fields: %v", err)
	}
	for i, user := range responseUsers {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("native SQL user at index %d unexpectedly exposes password", i)
		}
	}
}

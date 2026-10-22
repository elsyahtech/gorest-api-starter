package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersWithNativeCQLE2E(t *testing.T) {
	testServer, app := newScyllaUserTestServer(t)
	path := "/examples/native-queries/users/scylla"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty native CQL users response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty native CQL users data = %#v, want an empty array", emptyResult.Data)
	}

	createdAt := time.Now().UTC().Truncate(time.Millisecond)
	expectedUsers := []entity.User{
		{
			ID:           "user-native-cql-read-e2e-01",
			Email:        "native-cql-read-01@example.test",
			Password:     "hashed-test-password-01",
			DepartmentID: "department-native-cql-read-e2e-a",
			IsActive:     true,
			CreatedAt:    createdAt,
		},
		{
			ID:           "user-native-cql-read-e2e-02",
			Email:        "native-cql-read-02@example.test",
			Password:     "hashed-test-password-02",
			DepartmentID: "department-native-cql-read-e2e-b",
			IsActive:     false,
			CreatedAt:    createdAt.Add(time.Second),
		},
	}
	for _, user := range expectedUsers {
		if err := executeScyllaCQL(t, app, `
			INSERT INTO users (id, email, password, department_id, is_active, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`,
			user.ID,
			user.Email,
			user.Password,
			user.DepartmentID,
			user.IsActive,
			user.CreatedAt,
		); err != nil {
			t.Fatalf("seed ScyllaDB user %q: %v", user.ID, err)
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
		t.Fatalf("decode native CQL users response: %v", err)
	}
	if len(result.Data) != len(expectedUsers) {
		t.Fatalf("native CQL user count = %d, want %d; data=%+v",
			len(result.Data), len(expectedUsers), result.Data)
	}

	usersByID := make(map[string]entity.UserResponse, len(result.Data))
	for _, user := range result.Data {
		usersByID[user.ID] = user
	}
	for _, expected := range expectedUsers {
		user, found := usersByID[expected.ID]
		if !found {
			t.Errorf("native CQL response is missing user %q", expected.ID)
			continue
		}
		if user.Email != expected.Email ||
			user.DepartmentID != expected.DepartmentID ||
			user.IsActive != expected.IsActive ||
			user.CreatedAt == "" {
			t.Errorf("native CQL user response = %+v, does not match fixture %+v",
				user, expected)
		}
		if user.UpdatedAt != nil || user.DeletedAt != nil {
			t.Errorf("native CQL nullable timestamps should be nil; updatedAt=%v deletedAt=%v",
				user.UpdatedAt, user.DeletedAt)
		}
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode native CQL response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode native CQL user fields: %v", err)
	}
	for i, user := range responseUsers {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("native CQL user at index %d unexpectedly exposes password", i)
		}
	}
}

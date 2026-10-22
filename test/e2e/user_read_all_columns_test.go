package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersWithAllColumnsE2E(t *testing.T) {
	testServer, _ := newUserTestServer(t)
	path := "/examples/orm-queries/users/select-all"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty all-columns response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty all-columns data = %#v, want an empty array", emptyResult.Data)
	}

	expected := entity.User{
		ID:           "user-select-all-e2e-01",
		Email:        "select-all-e2e@example.test",
		Password:     "plain-text-select-all-password",
		DepartmentID: "department-select-all-e2e-01",
		IsActive:     true,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("encode create user request: %v", err)
	}
	createResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users",
		body,
	)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orm-queries/users status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
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
		t.Fatalf("decode all-columns response: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("all-columns user count = %d, want 1; body: %s",
			len(result.Data), response.Body)
	}
	user := result.Data[0]
	if user.ID != expected.ID ||
		user.Email != expected.Email ||
		user.DepartmentID != expected.DepartmentID ||
		user.IsActive != expected.IsActive {
		t.Fatalf("all-columns user response = %+v, does not match user", user)
	}
	if user.CreatedAt == "" {
		t.Fatal("all-columns response createdAt is empty")
	}
	if user.UpdatedAt != nil || user.DeletedAt != nil {
		t.Fatalf("new user's nullable timestamps should be null; updatedAt=%v deletedAt=%v",
			user.UpdatedAt, user.DeletedAt)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode all-columns response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode all-columns user fields: %v", err)
	}
	if _, hasPassword := responseUsers[0]["password"]; hasPassword {
		t.Fatal("all-columns response unexpectedly exposes password")
	}
	if _, hasCreatedAt := responseUsers[0]["createdAt"]; !hasCreatedAt {
		t.Fatal("all-columns response does not include createdAt")
	}
	if _, hasUpdatedAt := responseUsers[0]["updatedAt"]; !hasUpdatedAt {
		t.Fatal("all-columns response does not include updatedAt")
	}
	if _, hasDeletedAt := responseUsers[0]["deletedAt"]; !hasDeletedAt {
		t.Fatal("all-columns response does not include deletedAt")
	}
}

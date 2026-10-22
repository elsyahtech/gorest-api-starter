package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersWithSelectedColumnsE2E(t *testing.T) {
	testServer, _ := newUserTestServer(t)
	path := "/examples/orm-queries/users/select-columns"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty selected-columns response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty selected-columns data = %#v, want an empty array", emptyResult.Data)
	}

	expected := entity.User{
		ID:           "user-select-columns-e2e-01",
		Email:        "select-columns-e2e@example.test",
		Password:     "plain-text-select-columns-password",
		DepartmentID: "department-select-columns-e2e-01",
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
		t.Fatalf("decode selected-columns response: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("selected-columns user count = %d, want 1; body: %s",
			len(result.Data), response.Body)
	}
	user := result.Data[0]
	if user.ID != expected.ID || user.Email != expected.Email {
		t.Fatalf("selected-columns user = %+v, want ID %q and email %q",
			user, expected.ID, expected.Email)
	}
	if user.DepartmentID != "" || user.IsActive || user.CreatedAt != "" ||
		user.UpdatedAt != nil || user.DeletedAt != nil {
		t.Fatalf("unselected user fields should retain zero values; got %+v", user)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode selected-columns response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode selected-columns user fields: %v", err)
	}
	for _, field := range []string{"id", "email"} {
		if _, found := responseUsers[0][field]; !found {
			t.Errorf("selected-columns response does not include %q", field)
		}
	}
	if _, hasPassword := responseUsers[0]["password"]; hasPassword {
		t.Error("selected-columns response unexpectedly exposes password")
	}
}

package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersByLimitE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	path := "/examples/orm-queries/users/limit"

	assertLimitedUserResponse(t, sendRequest(t, testServer, http.MethodGet, path+"?limit=2", nil), 0)

	users := []entity.User{
		{ID: "user-limit-e2e-01", Email: "user-limit-01@example.test", Password: "test-password-01", DepartmentID: "department-limit-e2e-01", IsActive: true},
		{ID: "user-limit-e2e-02", Email: "user-limit-02@example.test", Password: "test-password-02", DepartmentID: "department-limit-e2e-02", IsActive: false},
		{ID: "user-limit-e2e-03", Email: "user-limit-03@example.test", Password: "test-password-03", DepartmentID: "department-limit-e2e-03", IsActive: true},
	}
	for _, user := range users {
		ensureE2EUserTimestamp(&user)
		query := app.Database().Table("users").Create(&user)
		if query.Error != nil {
			t.Fatalf("seed user %s: %s: %v", user.ID, query.Message, query.Error)
		}
	}

	response := sendRequest(t, testServer, http.MethodGet, path+"?limit=2", nil)
	assertLimitedUserResponse(t, response, 2)

	var result struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode limited users response: %v", err)
	}
	expectedIDs := make(map[string]struct{}, len(users))
	for _, user := range users {
		expectedIDs[user.ID] = struct{}{}
	}
	for _, user := range result.Data {
		if _, found := expectedIDs[user.ID]; !found {
			t.Errorf("unexpected user in limited response: %+v", user)
		}
		if user.Email == "" {
			t.Errorf("limited user response is missing email: %+v", user)
		}
	}

	zeroLimit := sendRequest(t, testServer, http.MethodGet, path+"?limit=0", nil)
	assertLimitedUserResponse(t, zeroLimit, 0)

	for _, query := range []string{"", "?limit=invalid"} {
		invalidResponse := sendRequest(t, testServer, http.MethodGet, path+query, nil)
		if invalidResponse.StatusCode != http.StatusBadRequest {
			t.Errorf("GET %s%s status = %d, want %d; body: %s",
				path, query, invalidResponse.StatusCode, http.StatusBadRequest, invalidResponse.Body)
		}
	}
}

func assertLimitedUserResponse(t *testing.T, response testHTTPResponse, expectedCount int) {
	t.Helper()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("limited users status = %d, want %d; body: %s",
			response.StatusCode, http.StatusOK, response.Body)
	}

	var result struct {
		Data  entity.UserResponseArray `json:"data"`
		Total int                      `json:"Total"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode limited users response: %v", err)
	}
	if result.Data == nil {
		t.Fatalf("limited users data = null, want an array; body: %s", response.Body)
	}
	if len(result.Data) != expectedCount {
		t.Fatalf("limited users count = %d, want %d; body: %s",
			len(result.Data), expectedCount, response.Body)
	}
	if result.Total != expectedCount {
		t.Fatalf("limited users Total = %d, want %d; body: %s",
			result.Total, expectedCount, response.Body)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode limited response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode limited user fields: %v", err)
	}
	for i, user := range responseUsers {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("limited user at index %d unexpectedly exposes password", i)
		}
	}
}

package e2e_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersByLimitOffsetE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	path := "/examples/orm-queries/users/pagination"

	users := []entity.User{
		{ID: "user-page-e2e-03", Email: "page-03@example.test", Password: "test-password-03", DepartmentID: "department-page-e2e-03", IsActive: true},
		{ID: "user-page-e2e-01", Email: "page-01@example.test", Password: "test-password-01", DepartmentID: "department-page-e2e-01", IsActive: true},
		{ID: "user-page-e2e-04", Email: "page-04@example.test", Password: "test-password-04", DepartmentID: "department-page-e2e-04", IsActive: false},
		{ID: "user-page-e2e-02", Email: "page-02@example.test", Password: "test-password-02", DepartmentID: "department-page-e2e-02", IsActive: false},
	}
	for _, user := range users {
		query := app.Database().Table("users").Create(&user)
		if query.Error != nil {
			t.Fatalf("seed user %s: %s: %v", user.ID, query.Message, query.Error)
		}
	}

	firstPage := sendRequest(
		t,
		testServer,
		http.MethodGet,
		path+"?limit=2&offset=0",
		nil,
	)
	assertPaginationResponse(t, firstPage, []string{"user-page-e2e-01", "user-page-e2e-02"})

	secondPage := sendRequest(
		t,
		testServer,
		http.MethodGet,
		path+"?limit=2&offset=2",
		nil,
	)
	assertPaginationResponse(t, secondPage, []string{"user-page-e2e-03", "user-page-e2e-04"})

	emptyPage := sendRequest(
		t,
		testServer,
		http.MethodGet,
		path+"?limit=2&offset=10",
		nil,
	)
	assertPaginationResponse(t, emptyPage, []string{})

	invalidQueries := []string{
		path + "?limit=2",
		path + "?offset=0",
	}
	for _, query := range invalidQueries {
		response := sendRequest(t, testServer, http.MethodGet, query, nil)
		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("GET %s status = %d, want %d; body: %s",
				query, response.StatusCode, http.StatusBadRequest, response.Body)
		}
	}
}

func assertPaginationResponse(t *testing.T, response testHTTPResponse, expectedIDs []string) {
	t.Helper()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("pagination response status = %d, want %d; body: %s",
			response.StatusCode, http.StatusOK, response.Body)
	}

	var result struct {
		Data  entity.UserResponseArray `json:"data"`
		Total int                      `json:"Total"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode pagination response: %v", err)
	}
	if result.Data == nil {
		t.Fatalf("pagination data = null, want an array; body: %s", response.Body)
	}
	actualIDs := make([]string, len(result.Data))
	for i, user := range result.Data {
		actualIDs[i] = user.ID
		if user.Email == "" {
			t.Errorf("paginated user %q has no email", user.ID)
		}
	}
	if !reflect.DeepEqual(actualIDs, expectedIDs) {
		t.Fatalf("paginated user IDs = %v, want %v", actualIDs, expectedIDs)
	}
	if result.Total != len(expectedIDs) {
		t.Fatalf("pagination Total = %d, want %d; body: %s",
			result.Total, len(expectedIDs), response.Body)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode pagination response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode paginated user fields: %v", err)
	}
	for i, user := range responseUsers {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("paginated user at index %d unexpectedly exposes password", i)
		}
	}
}

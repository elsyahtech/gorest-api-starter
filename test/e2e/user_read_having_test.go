package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetDepartmentsByUserCountE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	path := "/examples/orm-queries/users/having"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.GroupByResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty department-count response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty department-count data = %#v, want an empty array", emptyResult.Data)
	}

	users := []entity.User{
		{ID: "user-having-e2e-01", Email: "having-01@example.test", Password: "test-password-01", DepartmentID: "department-having-above-threshold"},
		{ID: "user-having-e2e-02", Email: "having-02@example.test", Password: "test-password-02", DepartmentID: "department-having-above-threshold"},
		{ID: "user-having-e2e-03", Email: "having-03@example.test", Password: "test-password-03", DepartmentID: "department-having-above-threshold"},
		{ID: "user-having-e2e-04", Email: "having-04@example.test", Password: "test-password-04", DepartmentID: "department-having-at-threshold"},
		{ID: "user-having-e2e-05", Email: "having-05@example.test", Password: "test-password-05", DepartmentID: "department-having-at-threshold"},
		{ID: "user-having-e2e-06", Email: "having-06@example.test", Password: "test-password-06", DepartmentID: "department-having-below-threshold"},
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
		Data entity.GroupByResponseArray `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode department-count response: %v", err)
	}
	expected := entity.GroupByResponseArray{
		{DepartmentID: "department-having-above-threshold", UserCount: 3},
	}
	if len(result.Data) != len(expected) {
		t.Fatalf("departments above user-count threshold = %+v, want %+v", result.Data, expected)
	}
	if result.Data[0] != expected[0] {
		t.Fatalf("department count result = %+v, want %+v", result.Data[0], expected[0])
	}
}

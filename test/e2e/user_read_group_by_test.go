package e2e_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUserCountsByDepartmentE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	path := "/examples/orm-queries/users/group-by"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.GroupByResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty department counts response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty department counts = %#v, want an empty array", emptyResult.Data)
	}

	users := []entity.User{
		{ID: "user-group-by-e2e-01", Email: "group-by-01@example.test", Password: "test-password-01", DepartmentID: "department-group-by-a"},
		{ID: "user-group-by-e2e-02", Email: "group-by-02@example.test", Password: "test-password-02", DepartmentID: "department-group-by-a"},
		{ID: "user-group-by-e2e-03", Email: "group-by-03@example.test", Password: "test-password-03", DepartmentID: "department-group-by-a"},
		{ID: "user-group-by-e2e-04", Email: "group-by-04@example.test", Password: "test-password-04", DepartmentID: "department-group-by-b"},
		{ID: "user-group-by-e2e-05", Email: "group-by-05@example.test", Password: "test-password-05", DepartmentID: "department-group-by-b"},
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
		Data entity.GroupByResponseArray `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode department counts response: %v", err)
	}
	expected := entity.GroupByResponseArray{
		{DepartmentID: "department-group-by-a", UserCount: 3},
		{DepartmentID: "department-group-by-b", UserCount: 2},
	}
	if len(result.Data) != len(expected) {
		t.Fatalf("department count groups = %d, want %d; data=%+v",
			len(result.Data), len(expected), result.Data)
	}

	actualByDepartment := make(map[string]int64, len(result.Data))
	for _, summary := range result.Data {
		actualByDepartment[summary.DepartmentID] = summary.UserCount
	}
	expectedByDepartment := make(map[string]int64, len(expected))
	for _, summary := range expected {
		expectedByDepartment[summary.DepartmentID] = summary.UserCount
	}
	if !reflect.DeepEqual(actualByDepartment, expectedByDepartment) {
		t.Fatalf("department user counts = %v, want %v", actualByDepartment, expectedByDepartment)
	}
}

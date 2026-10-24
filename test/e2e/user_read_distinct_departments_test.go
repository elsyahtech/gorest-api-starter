package e2e_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetUsersWithDistinctDepartmentIDsE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	path := "/examples/orm-queries/users/distinct"

	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty distinct-departments response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty distinct-departments data = %#v, want an empty array", emptyResult.Data)
	}

	users := []entity.User{
		{ID: "user-distinct-e2e-01", Email: "distinct-01@example.test", Password: "test-password-01", DepartmentID: "department-distinct-a"},
		{ID: "user-distinct-e2e-02", Email: "distinct-02@example.test", Password: "test-password-02", DepartmentID: "department-distinct-a"},
		{ID: "user-distinct-e2e-03", Email: "distinct-03@example.test", Password: "test-password-03", DepartmentID: "department-distinct-b"},
		{ID: "user-distinct-e2e-04", Email: "distinct-04@example.test", Password: "test-password-04", DepartmentID: "department-distinct-b"},
		{ID: "user-distinct-e2e-05", Email: "distinct-05@example.test", Password: "test-password-05", DepartmentID: "department-distinct-b"},
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
		t.Fatalf("decode distinct-departments response: %v", err)
	}

	actualDepartmentIDs := make([]string, 0, len(result.Data))
	for _, user := range result.Data {
		actualDepartmentIDs = append(actualDepartmentIDs, user.DepartmentID)
		if user.ID != "" || user.Email != "" || user.IsActive || user.CreatedAt != "" {
			t.Errorf("distinct department result should contain only the selected column; got %+v", user)
		}
	}
	expectedDepartmentIDs := []string{"department-distinct-a", "department-distinct-b"}
	if len(actualDepartmentIDs) != len(expectedDepartmentIDs) {
		t.Fatalf("distinct department IDs = %v, want %v", actualDepartmentIDs, expectedDepartmentIDs)
	}
	actualSet := make(map[string]struct{}, len(actualDepartmentIDs))
	for _, id := range actualDepartmentIDs {
		if _, duplicate := actualSet[id]; duplicate {
			t.Fatalf("department ID %q appeared more than once in distinct results", id)
		}
		actualSet[id] = struct{}{}
	}
	expectedSet := map[string]struct{}{
		"department-distinct-a": {},
		"department-distinct-b": {},
	}
	if !reflect.DeepEqual(actualSet, expectedSet) {
		t.Fatalf("distinct department IDs = %v, want %v", actualSet, expectedSet)
	}
}

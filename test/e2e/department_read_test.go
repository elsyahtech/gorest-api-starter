package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestGetDepartmentsE2E(t *testing.T) {
	testServer := newDepartmentTestServer(t)

	emptyResponse := sendRequest(t, testServer, http.MethodGet, "/examples/departments", nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/departments status = %d, want %d; body: %s",
			emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}

	var emptyResult struct {
		Data entity.DepartmentResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty department list response: %v", err)
	}
	if emptyResult.Data == nil {
		t.Fatal("empty department list data is null, want an empty array")
	}
	if len(emptyResult.Data) != 0 {
		t.Fatalf("department count in empty list = %d, want 0", len(emptyResult.Data))
	}

	expectedDepartments := []entity.Department{
		{
			ID:          "dept-read-e2e-01",
			Name:        "E2E Engineering",
			Description: "Engineering department",
			IsActive:    true,
		},
		{
			ID:          "dept-read-e2e-02",
			Name:        "E2E Design",
			Description: "Design department",
			IsActive:    false,
		},
	}
	for _, department := range expectedDepartments {
		body, err := json.Marshal(department)
		if err != nil {
			t.Fatalf("encode department %q: %v", department.ID, err)
		}

		createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/departments", body)
		if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
			t.Fatalf("POST /examples/departments for %q status = %d, want 2xx; body: %s",
				department.ID, createResponse.StatusCode, createResponse.Body)
		}
	}

	listResponse := sendRequest(t, testServer, http.MethodGet, "/examples/departments", nil)
	if listResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/departments status = %d, want %d; body: %s",
			listResponse.StatusCode, http.StatusOK, listResponse.Body)
	}

	var listResult struct {
		Data entity.DepartmentResponseArray `json:"data"`
	}
	if err := json.Unmarshal(listResponse.Body, &listResult); err != nil {
		t.Fatalf("decode department list response: %v", err)
	}
	if len(listResult.Data) != len(expectedDepartments) {
		t.Fatalf("department count = %d, want %d", len(listResult.Data), len(expectedDepartments))
	}

	departmentsByID := make(map[string]entity.DepartmentResponse, len(listResult.Data))
	for _, department := range listResult.Data {
		departmentsByID[department.ID] = department
	}
	for _, expected := range expectedDepartments {
		actual, ok := departmentsByID[expected.ID]
		if !ok {
			t.Errorf("department %q not found in response", expected.ID)
			continue
		}
		if actual.Name != expected.Name ||
			actual.Description != expected.Description ||
			actual.IsActive != expected.IsActive {
			t.Errorf("department %q = %+v, want name %q, description %q, isActive %t",
				expected.ID, actual, expected.Name, expected.Description, expected.IsActive)
		}
	}
}

func TestGetDepartmentByIDE2E(t *testing.T) {
	testServer := newDepartmentTestServer(t)
	expected := entity.Department{
		ID:          "dept-read-by-id-e2e-01",
		Name:        "E2E Platform",
		Description: "Platform department",
		IsActive:    true,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("encode department: %v", err)
	}

	createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/departments", body)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/departments status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/departments/"+expected.ID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/departments/%s status = %d, want %d; body: %s",
			expected.ID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}

	var result struct {
		Data entity.DepartmentResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &result); err != nil {
		t.Fatalf("decode department response: %v", err)
	}
	if result.Data.ID != expected.ID ||
		result.Data.Name != expected.Name ||
		result.Data.Description != expected.Description ||
		result.Data.IsActive != expected.IsActive {
		t.Fatalf("department response = %+v, want ID %q, name %q, description %q, isActive %t",
			result.Data, expected.ID, expected.Name, expected.Description, expected.IsActive)
	}
	if result.Data.CreatedAt == "" {
		t.Fatal("department createdAt is empty, want a formatted creation timestamp")
	}
}

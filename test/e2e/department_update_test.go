package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestUpdateDepartmentE2E(t *testing.T) {
	testServer := newDepartmentTestServer(t)
	departmentID := "dept-update-e2e-01"
	initial := entity.Department{
		ID:          departmentID,
		Name:        "Original Department",
		Description: "Original description",
		IsActive:    true,
	}
	createBody, err := json.Marshal(initial)
	if err != nil {
		t.Fatalf("encode create department request: %v", err)
	}

	createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/departments", createBody)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/departments status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	updated := entity.Department{
		ID:          departmentID,
		Name:        "Updated Department",
		Description: "Updated description",
		IsActive:    false,
	}
	updateBody, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode update department request: %v", err)
	}

	updateResponse := sendRequest(
		t,
		testServer,
		http.MethodPut,
		"/examples/departments/"+departmentID,
		updateBody,
	)
	if updateResponse.StatusCode < http.StatusOK || updateResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PUT /examples/departments/%s status = %d, want 2xx; body: %s",
			departmentID, updateResponse.StatusCode, updateResponse.Body)
	}

	var updateResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(updateResponse.Body, &updateResult); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if !updateResult.Success {
		t.Fatal("update response success = false, want true")
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/departments/"+departmentID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/departments/%s status = %d, want %d; body: %s",
			departmentID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}

	var getResult struct {
		Data entity.DepartmentResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &getResult); err != nil {
		t.Fatalf("decode department response after update: %v", err)
	}
	if getResult.Data.ID != updated.ID ||
		getResult.Data.Name != updated.Name ||
		getResult.Data.Description != updated.Description ||
		getResult.Data.IsActive != updated.IsActive {
		t.Fatalf("department after update = %+v, want ID %q, name %q, description %q, isActive %t",
			getResult.Data, updated.ID, updated.Name, updated.Description, updated.IsActive)
	}
}

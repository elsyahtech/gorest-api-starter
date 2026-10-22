package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestPatchDepartmentE2E(t *testing.T) {
	testServer := newDepartmentTestServer(t)
	departmentID := "dept-patch-e2e-01"
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

	inactive := false
	patchBody, err := json.Marshal(entity.DepartmentPatch{IsActive: &inactive})
	if err != nil {
		t.Fatalf("encode patch department request: %v", err)
	}

	patchResponse := sendRequest(
		t,
		testServer,
		http.MethodPatch,
		"/examples/departments/"+departmentID,
		patchBody,
	)
	if patchResponse.StatusCode < http.StatusOK || patchResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PATCH /examples/departments/%s status = %d, want 2xx; body: %s",
			departmentID, patchResponse.StatusCode, patchResponse.Body)
	}

	var patchResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(patchResponse.Body, &patchResult); err != nil {
		t.Fatalf("decode patch response: %v", err)
	}
	if !patchResult.Success {
		t.Fatal("patch response success = false, want true")
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
		t.Fatalf("decode department response after patch: %v", err)
	}
	if getResult.Data.ID != initial.ID ||
		getResult.Data.Name != initial.Name ||
		getResult.Data.Description != initial.Description ||
		getResult.Data.IsActive {
		t.Fatalf("department after patch = %+v, want original fields retained and isActive=false",
			getResult.Data)
	}
	if getResult.Data.UpdatedAt == nil || *getResult.Data.UpdatedAt == "" {
		t.Fatal("department updatedAt is empty, want a formatted update timestamp")
	}

	emptyPatchResponse := sendRequest(
		t,
		testServer,
		http.MethodPatch,
		"/examples/departments/"+departmentID,
		[]byte(`{}`),
	)
	if emptyPatchResponse.StatusCode != http.StatusBadRequest {
		t.Fatalf("PATCH with empty object status = %d, want %d; body: %s",
			emptyPatchResponse.StatusCode, http.StatusBadRequest, emptyPatchResponse.Body)
	}
}

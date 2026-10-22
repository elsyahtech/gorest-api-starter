package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestCreateDepartmentE2E(t *testing.T) {
	testServer := newDepartmentTestServer(t)

	payload := entity.Department{
		ID:          "dept-e2e-01",
		Name:        "E2E Engineering",
		Description: "Created by the department endpoint E2E test",
		IsActive:    true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("encode create department request: %v", err)
	}

	postResponse := sendRequest(t, testServer, http.MethodPost, "/examples/departments", body)
	if postResponse.StatusCode < http.StatusOK || postResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/departments status = %d, want 2xx; body: %s",
			postResponse.StatusCode, postResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(postResponse.Body, &createResult); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("create response success = false, want true")
	}

	duplicatePayload := payload
	duplicatePayload.Name = "Must Not Replace Existing Department"
	duplicateBody, err := json.Marshal(duplicatePayload)
	if err != nil {
		t.Fatalf("encode duplicate department request: %v", err)
	}

	duplicateResponse := sendRequest(t, testServer, http.MethodPost, "/examples/departments", duplicateBody)
	if duplicateResponse.StatusCode < http.StatusBadRequest {
		t.Fatalf("duplicate POST status = %d, want 4xx or 5xx; body: %s",
			duplicateResponse.StatusCode, duplicateResponse.Body)
	}

	getResponse := sendRequest(t, testServer, http.MethodGet, "/examples/departments", nil)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/departments status = %d, want %d; body: %s",
			getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}

	var listResult struct {
		Data entity.DepartmentResponseArray `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &listResult); err != nil {
		t.Fatalf("decode department list response: %v", err)
	}
	if len(listResult.Data) != 1 {
		t.Fatalf("department count after duplicate create = %d, want 1", len(listResult.Data))
	}

	department := listResult.Data[0]
	if department.ID != payload.ID ||
		department.Name != payload.Name ||
		department.Description != payload.Description ||
		!department.IsActive {
		t.Fatalf("persisted department = %+v, does not match original request", department)
	}
}

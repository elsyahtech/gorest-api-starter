package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestDeleteDepartmentE2E(t *testing.T) {
	testServer := newDepartmentTestServer(t)
	payload := entity.Department{
		ID:          "dept-delete-e2e-01",
		Name:        "E2E Department to Delete",
		Description: "Created by the department delete E2E test",
		IsActive:    true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("encode create department request: %v", err)
	}

	createResponse := sendRequest(t, testServer, http.MethodPost, "/examples/departments", body)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/departments status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	deleteResponse := sendRequest(
		t,
		testServer,
		http.MethodDelete,
		"/examples/departments/"+payload.ID,
		nil,
	)
	if deleteResponse.StatusCode < http.StatusOK || deleteResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("DELETE /examples/departments/%s status = %d, want 2xx; body: %s",
			payload.ID, deleteResponse.StatusCode, deleteResponse.Body)
	}

	var deleteResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(deleteResponse.Body, &deleteResult); err != nil {
		t.Fatalf("decode delete response: %v", err)
	}
	if !deleteResult.Success {
		t.Fatal("delete response success = false, want true")
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
	if len(listResult.Data) != 0 {
		t.Fatalf("department count after delete = %d, want 0", len(listResult.Data))
	}
}

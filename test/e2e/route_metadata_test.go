package e2e_test

import (
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestDepartmentHeadRoutesE2E(t *testing.T) {
	testServer := newDepartmentTestServer(t)
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodHead, "/examples/departments", nil),
		http.StatusOK, "")
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodHead, "/examples/departments/missing", nil),
		http.StatusNotFound, "")
}

func TestDepartmentOptionsRoutesE2E(t *testing.T) {
	testServer := newDepartmentTestServer(t)
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodOptions, "/examples/departments", nil),
		http.StatusNoContent, "GET, HEAD, POST, OPTIONS")
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodOptions, "/examples/departments/missing", nil),
		http.StatusNoContent, "GET, HEAD, PUT, PATCH, DELETE, OPTIONS")
}

func TestUserHeadRoutesE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	user := entity.User{
		ID:           "user-head-e2e-01",
		Email:        "user-head@example.test",
		Password:     "test-password",
		DepartmentID: "department-head-e2e-01",
		IsActive:     true,
	}
	query := app.Database().Table("users").Create(&user)
	if query.Error != nil {
		t.Fatalf("seed user %s: %s: %v", user.ID, query.Message, query.Error)
	}

	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodHead, "/examples/orm-queries/users", nil),
		http.StatusOK, "")
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodHead, "/examples/orm-queries/users/"+user.ID, nil),
		http.StatusOK, "")
}

func TestUserOptionsRoutesE2E(t *testing.T) {
	testServer, _ := newUserTestServer(t)
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodOptions, "/examples/orm-queries/users", nil),
		http.StatusNoContent, "GET, HEAD, POST, OPTIONS")
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodOptions, "/examples/orm-queries/users/missing", nil),
		http.StatusNoContent, "GET, HEAD, PUT, PATCH, DELETE, OPTIONS")
}

func TestOrderHeadRoutesE2E(t *testing.T) {
	testServer, userID := newOrderTestServer(t)
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodHead, "/examples/orders", nil),
		http.StatusOK, "")
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodHead, "/examples/users/"+userID+"/orders", nil),
		http.StatusOK, "")
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodHead, "/examples/orders/missing", nil),
		http.StatusNotFound, "")
}

func TestOrderOptionsRoutesE2E(t *testing.T) {
	testServer, userID := newOrderTestServer(t)
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodOptions, "/examples/orders", nil),
		http.StatusNoContent, "GET, HEAD, POST, OPTIONS")
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodOptions, "/examples/users/"+userID+"/orders", nil),
		http.StatusNoContent, "GET, HEAD, OPTIONS")
	assertMetadataRequest(t, sendRequest(t, testServer, http.MethodOptions, "/examples/orders/missing", nil),
		http.StatusNoContent, "GET, HEAD, PUT, PATCH, DELETE, OPTIONS")
}

func assertMetadataRequest(t *testing.T, response testHTTPResponse, expectedStatus int, expectedAllow string) {
	t.Helper()
	if response.StatusCode != expectedStatus {
		t.Fatalf("metadata response status = %d, want %d; body: %s",
			response.StatusCode, expectedStatus, response.Body)
	}
	if len(response.Body) != 0 {
		t.Errorf("metadata response body = %q, want no body", response.Body)
	}
	if expectedAllow != "" && response.Header.Get("Allow") != expectedAllow {
		t.Errorf("metadata Allow header = %q, want %q", response.Header.Get("Allow"), expectedAllow)
	}
}

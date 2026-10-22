package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
)

func TestDeleteUserWithNativeCQLE2E(t *testing.T) {
	testServer, app := newScyllaUserTestServer(t)
	userToDelete := "user-delete-native-cql-e2e-01"
	userToKeep := "user-delete-native-cql-e2e-02"
	for _, userID := range []string{userToDelete, userToKeep} {
		if err := executeScyllaCQL(t, app, `
			INSERT INTO users (id, email, password, department_id, is_active, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`,
			userID,
			userID+"@example.test",
			"hashed-test-password",
			"department-delete-native-cql-e2e",
			true,
			time.Now(),
		); err != nil {
			t.Fatalf("seed ScyllaDB user %q: %v", userID, err)
		}
	}

	deleteResponse := sendRequest(
		t,
		testServer,
		http.MethodDelete,
		"/examples/native-queries/users/scylla/"+userToDelete,
		nil,
	)
	if deleteResponse.StatusCode < http.StatusOK || deleteResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("DELETE /examples/native-queries/users/scylla/%s status = %d, want 2xx; body: %s",
			userToDelete, deleteResponse.StatusCode, deleteResponse.Body)
	}

	var deleteResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(deleteResponse.Body, &deleteResult); err != nil {
		t.Fatalf("decode native CQL delete response: %v", err)
	}
	if !deleteResult.Success {
		t.Fatal("native CQL delete response success = false, want true")
	}

	if scyllaUserExists(t, app, userToDelete) {
		t.Fatalf("ScyllaDB user %q still exists after delete", userToDelete)
	}
	if !scyllaUserExists(t, app, userToKeep) {
		t.Fatalf("unaffected ScyllaDB user %q was deleted", userToKeep)
	}

	repeatedDeleteResponse := sendRequest(
		t,
		testServer,
		http.MethodDelete,
		"/examples/native-queries/users/scylla/"+userToDelete,
		nil,
	)
	if repeatedDeleteResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("repeated DELETE status = %d, want %d; body: %s",
			repeatedDeleteResponse.StatusCode, http.StatusNotFound, repeatedDeleteResponse.Body)
	}
}

func scyllaUserExists(t *testing.T, app *gorest.App, userID string) bool {
	t.Helper()

	iter, _, _, err := app.Database().QueryCQL(t.Context(),
		"SELECT id FROM users WHERE id = ?", userID,
	)
	if err != nil {
		t.Fatalf("query ScyllaDB user %q: %v", userID, err)
	}

	var foundID string
	found := iter.Scan(&foundID)
	if err := iter.Close(); err != nil {
		t.Fatalf("close ScyllaDB user query for %q: %v", userID, err)
	}
	return found
}

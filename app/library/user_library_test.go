package library

import (
	"net/http"
	"testing"
)

func TestCreateUserRejectsNilPayload(t *testing.T) {
	lib := NewUserLibrary(nil, nil)

	message, status, err := lib.CreateUser(nil)
	if err == nil || status != http.StatusBadRequest || message == "" {
		t.Fatalf("CreateUser(nil) = (%q, %d, %v), want a bad-request error", message, status, err)
	}
}

func TestLoginRejectsNilPayload(t *testing.T) {
	lib := NewUserLibrary(nil, nil)

	message, status, err := lib.Login(nil, nil)
	if err == nil || status != http.StatusBadRequest || message == "" {
		t.Fatalf("Login(nil) = (%q, %d, %v), want a bad-request error", message, status, err)
	}
}

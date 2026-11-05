package model

import (
	"net/http"
	"testing"
)

func TestCreateUserRejectsNilPayload(t *testing.T) {
	model := &UserModel{}

	message, status, err := model.CreateUser(nil)
	if err == nil || status != http.StatusBadRequest || message == "" {
		t.Fatalf("CreateUser(nil) = (%q, %d, %v), want a bad-request error", message, status, err)
	}
}

func TestGetUserByEmailRejectsEmptyEmail(t *testing.T) {
	model := &UserModel{}

	user, message, status, err := model.GetUserByEmail("")
	if err == nil || user != nil || status != http.StatusBadRequest || message == "" {
		t.Fatalf("GetUserByEmail(empty) = (%v, %q, %d, %v), want a bad-request error", user, message, status, err)
	}
}

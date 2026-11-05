package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/server"
)

type userUseCaseStub struct {
	user    *entity.User
	users   *[]entity.User
	message string
	status  int
	err     error
}

func (stub *userUseCaseStub) Login(_ *gorest.Context, user *entity.User) (string, int, error) {
	stub.user = user
	return stub.message, stub.status, stub.err
}

func (stub *userUseCaseStub) CreateUser(user *entity.User) (string, int, error) {
	stub.user = user
	return stub.message, stub.status, stub.err
}

func (stub *userUseCaseStub) GetUsers() (*[]entity.User, string, int, error) {
	return stub.users, stub.message, stub.status, stub.err
}

func TestCreateUserDecodesRequestAndReturnsResponse(t *testing.T) {
	useCase := &userUseCaseStub{message: "User created", status: http.StatusCreated}
	handler := NewUserHandler(nil, useCase)
	app, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("server.New() error = %v", err)
	}
	app.Post("/users/create", handler.CreateUser)

	request := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader(
		`{"id":"user-1","email":"user@example.test","password":"password","isActive":true}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.SRV.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated || useCase.user == nil ||
		useCase.user.Email != "user@example.test" || useCase.user.Password != "password" {
		t.Fatalf("CreateUser() response = %d, decoded user = %+v", response.StatusCode, useCase.user)
	}

	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if body["success"] != true || body["message"] != "User created" {
		t.Errorf("response body = %v, want successful create response", body)
	}
}

func TestGetUsersReturnsPublicUserFields(t *testing.T) {
	users := []entity.User{{ID: "user-1", Email: "user@example.test", Password: "private"}}
	useCase := &userUseCaseStub{users: &users, message: "Users retrieved", status: http.StatusOK}
	handler := NewUserHandler(nil, useCase)
	app, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("server.New() error = %v", err)
	}
	app.Get("/users/", handler.GetUsers)

	response, err := app.SRV.Test(httptest.NewRequest(http.MethodGet, "/users/", nil))
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /users/ status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal response body: %v", err)
	}
	if strings.Contains(string(encoded), "private") || body["success"] != true {
		t.Errorf("response body exposed private fields or failed: %s", encoded)
	}
}

func TestLoginDecodesRequestAndReturnsSuccess(t *testing.T) {
	useCase := &userUseCaseStub{message: "Login successful", status: http.StatusOK}
	handler := NewUserHandler(nil, useCase)
	app, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("server.New() error = %v", err)
	}
	app.Post("/users/login", handler.Login)

	request := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(
		`{"email":"user@example.test","password":"password"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.SRV.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK || useCase.user == nil ||
		useCase.user.Email != "user@example.test" || useCase.user.Password != "password" {
		t.Fatalf("Login() response = %d, decoded user = %+v", response.StatusCode, useCase.user)
	}
}

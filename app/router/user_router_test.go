package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest/filter"
	"github.com/elsyahtech/gorest/server"
)

type routeUserUseCaseStub struct{}

func (routeUserUseCaseStub) Login(*gorest.Context, *entity.User) (string, int, error) {
	return "Login successful", http.StatusOK, nil
}

func (routeUserUseCaseStub) CreateUser(*entity.User) (string, int, error) {
	return "User created", http.StatusOK, nil
}

func (routeUserUseCaseStub) GetUsers() (*[]entity.User, string, int, error) {
	users := []entity.User{}
	return &users, "Users retrieved", http.StatusOK, nil
}

func TestRegisterUserRoutesProtectsListAndRegistersPublicPostRoutes(t *testing.T) {
	app := gorest.New(gorest.ConfigDefault)
	signatureKey := "unit-test-key-that-is-not-used-for-real-authentication"
	filterConfig := filter.ConfigDefault
	filterConfig.SignatureKey = &signatureKey
	app.Use(filter.Run(filterConfig))
	t.Cleanup(func() {
		if _, _, err := app.Close(); err != nil {
			t.Errorf("close GOREST application: %v", err)
		}
	})

	httpServer, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("server.New() error = %v", err)
	}
	userHandler := handler.NewUserHandler(app, routeUserUseCaseStub{})
	RegisterUserRoutes(app, httpServer.Group("/users"), userHandler)

	response, err := httpServer.SRV.Test(httptest.NewRequest(http.MethodGet, "/users/", nil))
	if err != nil {
		t.Fatalf("GET /users/ error = %v", err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated GET /users/ status = %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}

	for _, route := range []struct {
		path string
		body string
	}{
		{path: "/users/create", body: `{"id":"user-1","email":"user@example.test","password":"password"}`},
		{path: "/users/login", body: `{"email":"user@example.test","password":"password"}`},
	} {
		request := httptest.NewRequest(http.MethodPost, route.path, strings.NewReader(route.body))
		request.Header.Set("Content-Type", "application/json")
		response, err := httpServer.SRV.Test(request)
		if err != nil {
			t.Fatalf("POST %s error = %v", route.path, err)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusOK {
			t.Errorf("POST %s status = %d, want %d", route.path, response.StatusCode, http.StatusOK)
		}
	}
}

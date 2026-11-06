package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest/filter"
	"github.com/elsyahtech/gorest/server"
)

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
	userHandler := handler.NewUserHandler(app, nil)
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
	}{
		{path: "/users/create"},
		{path: "/users/login"},
	} {
		request := httptest.NewRequest(http.MethodPost, route.path, strings.NewReader("{"))
		request.Header.Set("Content-Type", "application/json")
		response, err := httpServer.SRV.Test(request)
		if err != nil {
			t.Fatalf("POST %s error = %v", route.path, err)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("POST %s status = %d, want an unauthenticated request to reach the handler and return %d",
				route.path, response.StatusCode, http.StatusBadRequest)
		}
	}
}

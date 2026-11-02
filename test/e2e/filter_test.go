package e2e_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/filter"
	"github.com/elsyahtech/gorest/helper"
	"github.com/elsyahtech/gorest/server"
)

func TestAuthenticationFilterRequiresAndAcceptsAccessToken(t *testing.T) {
	signatureKey := "gorest-api-starter-filter-integration-test-key"
	filterConfig := filter.ConfigDefault
	filterConfig.SignatureKey = &signatureKey
	filterConfig.Issuer = "filter-integration-test"
	filterConfig.CookieSecure = false
	filterConfig.LoggerConfig.DisableLogging = true

	app := gorest.New(gorest.Config{AppName: "filter-integration-test"})
	app.Use(filterConfig)
	t.Cleanup(func() {
		if message, _, err := app.Close(); err != nil {
			t.Errorf("close GOREST app (%s): %v", message, err)
		}
	})

	httpServer, _, err := server.New(nil)
	if err != nil {
		t.Fatalf("create HTTP test server: %v", err)
	}
	httpServer.Post("/token", func(ctx *gorest.Context) error {
		message, status, err := app.Filter().CreateToken(*ctx, "integration-test-user")
		if err != nil {
			return ctx.Status(status).SendString(message)
		}
		return ctx.Status(status).SendString(message)
	})
	httpServer.Get("/protected", app.Filter().Authentication(), func(ctx *gorest.Context) error {
		return ctx.SendStatus(http.StatusNoContent)
	})

	unauthenticatedResponse, err := httpServer.SRV.Test(httptest.NewRequest(http.MethodGet, "/protected", nil))
	if err != nil {
		t.Fatalf("request protected route without token: %v", err)
	}
	defer unauthenticatedResponse.Body.Close()
	if unauthenticatedResponse.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated response status = %d, want %d",
			unauthenticatedResponse.StatusCode, http.StatusUnauthorized)
	}

	invalidTokenRequest := httptest.NewRequest(http.MethodGet, "/protected", nil)
	invalidTokenRequest.AddCookie(&http.Cookie{Name: "access", Value: "invalid-token"})
	invalidTokenResponse, err := httpServer.SRV.Test(invalidTokenRequest)
	if err != nil {
		t.Fatalf("request protected route with invalid token: %v", err)
	}
	defer invalidTokenResponse.Body.Close()
	if invalidTokenResponse.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid-token response status = %d, want %d",
			invalidTokenResponse.StatusCode, http.StatusUnauthorized)
	}

	tokenResponse, err := httpServer.SRV.Test(httptest.NewRequest(http.MethodPost, "/token", nil))
	if err != nil {
		t.Fatalf("request token endpoint: %v", err)
	}
	defer tokenResponse.Body.Close()
	if tokenResponse.StatusCode != http.StatusOK {
		t.Fatalf("token response status = %d, want %d", tokenResponse.StatusCode, http.StatusOK)
	}

	var accessCookie *http.Cookie
	for _, cookie := range tokenResponse.Cookies() {
		if cookie.Name == "access" {
			accessCookie = cookie
			break
		}
	}
	if accessCookie == nil || accessCookie.Value == "" {
		t.Fatalf("token response did not contain an access cookie: %v", tokenResponse.Header.Values("Set-Cookie"))
	}

	authenticatedRequest := httptest.NewRequest(http.MethodGet, "/protected", nil)
	authenticatedRequest.AddCookie(accessCookie)
	authenticatedResponse, err := httpServer.SRV.Test(authenticatedRequest)
	if err != nil {
		t.Fatalf("request protected route with access token: %v", err)
	}
	defer authenticatedResponse.Body.Close()
	if authenticatedResponse.StatusCode != http.StatusNoContent {
		t.Fatalf("authenticated response status = %d, want %d",
			authenticatedResponse.StatusCode, http.StatusNoContent)
	}
}

func TestUserLoginAndProtectedRouteE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)

	const password = "login-e2e-password"
	passwordHash, message, status, err := helper.GenerateArgon(password)
	if err != nil {
		t.Fatalf("hash test user password: %s (%d): %v", message, status, err)
	}
	user := entity.User{
		ID:           "user-login-e2e",
		Email:        "user-login@example.test",
		Password:     passwordHash,
		DepartmentID: "department-login-e2e",
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
	}
	query := app.Database().Table("users").Create(&user)
	if query.Error != nil {
		t.Fatalf("seed login user: %s: %v", query.Message, query.Error)
	}

	const protectedPath = "/examples/orm-queries/users/with-authentication"
	unauthenticated := sendRequest(t, testServer, http.MethodGet, protectedPath, nil)
	if unauthenticated.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated protected route status = %d, want %d; body: %s",
			unauthenticated.StatusCode, http.StatusUnauthorized, unauthenticated.Body)
	}

	loginBody, err := json.Marshal(map[string]string{
		"email":    user.Email,
		"password": password,
	})
	if err != nil {
		t.Fatalf("encode login request: %v", err)
	}
	invalidLoginBody, err := json.Marshal(map[string]string{
		"email":    user.Email,
		"password": "incorrect-password",
	})
	if err != nil {
		t.Fatalf("encode invalid login request: %v", err)
	}
	invalidLogin := sendRequest(t, testServer, http.MethodPost, "/examples/login", invalidLoginBody)
	if invalidLogin.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid login status = %d, want %d; body: %s",
			invalidLogin.StatusCode, http.StatusUnauthorized, invalidLogin.Body)
	}

	loginResponse := sendRequest(t, testServer, http.MethodPost, "/examples/login", loginBody)
	if loginResponse.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, want %d; body: %s",
			loginResponse.StatusCode, http.StatusOK, loginResponse.Body)
	}

	var accessCookie *http.Cookie
	for _, cookie := range (&http.Response{Header: loginResponse.Header}).Cookies() {
		if cookie.Name == "access" {
			accessCookie = cookie
			break
		}
	}
	if accessCookie == nil || accessCookie.Value == "" {
		t.Fatalf("login response did not contain an access cookie: %v",
			loginResponse.Header.Values("Set-Cookie"))
	}

	authenticatedRequest := httptest.NewRequest(http.MethodGet, protectedPath, nil)
	authenticatedRequest.AddCookie(accessCookie)
	authenticatedResponse, err := testServer.SRV.Test(authenticatedRequest)
	if err != nil {
		t.Fatalf("request protected route with login token: %v", err)
	}
	defer func() {
		if err := authenticatedResponse.Body.Close(); err != nil {
			t.Errorf("close authenticated response body: %v", err)
		}
	}()
	if authenticatedResponse.StatusCode != http.StatusOK {
		t.Fatalf("authenticated protected route status = %d, want %d",
			authenticatedResponse.StatusCode, http.StatusOK)
	}
}

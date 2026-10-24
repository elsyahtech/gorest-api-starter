package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest-api-starter/app/router"
	gorestlog "github.com/elsyahtech/gorest/log"
	gorestredis "github.com/elsyahtech/gorest/redis"
	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func TestRedisUserCRUDE2E(t *testing.T) {
	if os.Getenv("GOREST_REDIS_E2E") != "1" {
		t.Skip("set GOREST_REDIS_E2E=1 to run the Redis-backed E2E test")
	}

	testServer, app := newRedisUserTestServer(t)
	redisClient, _, _, err := app.Redis().QueryRedis()
	if err != nil {
		t.Fatalf("get Redis client: %v", err)
	}

	userID := fmt.Sprintf("redis-e2e-%d", time.Now().UnixNano())
	redisKey := "user:" + userID
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := redisClient.Del(ctx, redisKey).Err(); err != nil {
			t.Errorf("delete Redis E2E fixture: %v", err)
		}
	})

	initial := entity.User{
		ID:           userID,
		Email:        userID + "@example.test",
		Password:     "initial-redis-password",
		DepartmentID: "dept-redis-before",
		IsActive:     true,
	}
	createResponse := sendRedisUserRequest(t, testServer, http.MethodPost, "/examples/redis-queries/users", initial)
	if createResponse.StatusCode != http.StatusOK {
		t.Fatalf("create user status = %d, want %d; body: %s", createResponse.StatusCode, http.StatusOK, createResponse.Body)
	}

	duplicate := initial
	duplicate.Email = "duplicate-redis@example.test"
	duplicateResponse := sendRedisUserRequest(t, testServer, http.MethodPost, "/examples/redis-queries/users", duplicate)
	if duplicateResponse.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate create status = %d, want %d; body: %s", duplicateResponse.StatusCode, http.StatusConflict, duplicateResponse.Body)
	}

	initialResponse := sendRequest(t, testServer, http.MethodGet, "/examples/redis-queries/users/"+userID, nil)
	initialUser := decodeRedisUserResponse(t, initialResponse, http.StatusOK)
	if initialUser.Email != initial.Email || initialUser.DepartmentID != initial.DepartmentID {
		t.Fatalf("duplicate create overwrote the existing user: %+v", initialUser)
	}

	missingID := userID + "-missing"
	missingUpdate := initial
	missingUpdate.ID = missingID
	missingResponse := sendRedisUserRequest(t, testServer, http.MethodPut, "/examples/redis-queries/users/"+missingID, missingUpdate)
	if missingResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("update missing user status = %d, want %d; body: %s", missingResponse.StatusCode, http.StatusNotFound, missingResponse.Body)
	}

	mismatchedID := initial
	mismatchedID.ID = userID + "-other"
	mismatchResponse := sendRedisUserRequest(t, testServer, http.MethodPut, "/examples/redis-queries/users/"+userID, mismatchedID)
	if mismatchResponse.StatusCode != http.StatusBadRequest {
		t.Fatalf("update with mismatched ID status = %d, want %d; body: %s", mismatchResponse.StatusCode, http.StatusBadRequest, mismatchResponse.Body)
	}

	updated := initial
	updated.Email = "updated-" + initial.Email
	updated.Password = "updated-redis-password"
	updated.DepartmentID = "dept-redis-after"
	updated.IsActive = false
	updated.CreatedAt = time.Time{}
	updateResponse := sendRedisUserRequest(t, testServer, http.MethodPut, "/examples/redis-queries/users/"+userID, updated)
	if updateResponse.StatusCode != http.StatusOK {
		t.Fatalf("update user status = %d, want %d; body: %s", updateResponse.StatusCode, http.StatusOK, updateResponse.Body)
	}

	updatedResponse := sendRequest(t, testServer, http.MethodGet, "/examples/redis-queries/users/"+userID, nil)
	updatedUser := decodeRedisUserResponse(t, updatedResponse, http.StatusOK)
	if updatedUser.ID != userID || updatedUser.Email != updated.Email || updatedUser.DepartmentID != updated.DepartmentID || updatedUser.IsActive {
		t.Fatalf("updated user = %+v, want updated values and the URL ID", updatedUser)
	}
	if updatedUser.CreatedAt != initialUser.CreatedAt {
		t.Fatalf("updated createdAt = %q, want original %q", updatedUser.CreatedAt, initialUser.CreatedAt)
	}

	listResponse := sendRequest(t, testServer, http.MethodGet, "/examples/redis-queries/users", nil)
	if listResponse.StatusCode != http.StatusOK {
		t.Fatalf("list users status = %d, want %d; body: %s", listResponse.StatusCode, http.StatusOK, listResponse.Body)
	}
	var listResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(listResponse.Body, &listResult); err != nil {
		t.Fatalf("decode Redis user list: %v", err)
	}
	found := false
	for _, user := range listResult.Data {
		if user.ID == userID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Redis user list does not contain %q", userID)
	}

	deleteResponse := sendRequest(t, testServer, http.MethodDelete, "/examples/redis-queries/users/"+userID, nil)
	if deleteResponse.StatusCode != http.StatusOK {
		t.Fatalf("delete user status = %d, want %d; body: %s", deleteResponse.StatusCode, http.StatusOK, deleteResponse.Body)
	}
	deletedResponse := sendRequest(t, testServer, http.MethodGet, "/examples/redis-queries/users/"+userID, nil)
	if deletedResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("get deleted user status = %d, want %d; body: %s", deletedResponse.StatusCode, http.StatusNotFound, deletedResponse.Body)
	}
	secondDelete := sendRequest(t, testServer, http.MethodDelete, "/examples/redis-queries/users/"+userID, nil)
	if secondDelete.StatusCode != http.StatusNotFound {
		t.Fatalf("delete missing user status = %d, want %d; body: %s", secondDelete.StatusCode, http.StatusNotFound, secondDelete.Body)
	}
}

func newRedisUserTestServer(t *testing.T) (*server.Server, *gorest.App) {
	t.Helper()

	host := os.Getenv("GOREST_REDIS_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("GOREST_REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	logConfig := gorestlog.ConfigDefault
	logConfig.LogDir = filepath.Join(t.TempDir(), "logs")
	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter-user-redis-e2e",
		Version:     "test",
		Environment: gorest.Testing,
		Timezone:    "UTC",
	})
	app.Use(
		gorestlog.Run(logConfig),
		gorestredis.Run(gorestredis.Config{
			Host:    host,
			Port:    port,
			Disable: false,
			Timeout: 5 * time.Second,
		}),
	)
	t.Cleanup(func() {
		if message, _, err := app.Close(); err != nil {
			t.Errorf("close Redis E2E application (%s): %v", message, err)
		}
	})

	serverConfig := server.ConfigDefault
	limiterConfig := limiter.ConfigDefault
	limiterConfig.Next = func(fiber.Ctx) bool { return true }
	serverConfig.Limiter = &limiterConfig
	testServer, _, err := server.New(&serverConfig)
	if err != nil {
		t.Fatalf("create Redis E2E HTTP server: %v", err)
	}
	userModel := model.NewUserModel(app)
	userLibrary := library.NewUserLibrary(userModel)
	userHandler := handler.NewUserHandler(app, userLibrary)
	router.RegisterUserRoutes(testServer.Group("/examples"), userHandler)
	return testServer, app
}

func sendRedisUserRequest(t *testing.T, testServer *server.Server, method, path string, user entity.User) testHTTPResponse {
	t.Helper()
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("encode Redis user request: %v", err)
	}
	return sendRequest(t, testServer, method, path, body)
}

func decodeRedisUserResponse(t *testing.T, response testHTTPResponse, wantStatus int) entity.UserResponse {
	t.Helper()
	if response.StatusCode != wantStatus {
		t.Fatalf("Redis user response status = %d, want %d; body: %s", response.StatusCode, wantStatus, response.Body)
	}
	var result struct {
		Data entity.UserResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode Redis user response: %v", err)
	}
	return result.Data
}

package e2e_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest-api-starter/app/router"
	"github.com/elsyahtech/gorest/database"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/server"
)

func TestCreateUserWithTransactionE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	expected := entity.User{
		ID:           "user-transaction-e2e-01",
		Email:        "transaction-e2e@example.test",
		Password:     "plain-text-test-password",
		DepartmentID: "department-transaction-e2e-01",
		IsActive:     true,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("encode create user request: %v", err)
	}

	createResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users/transaction",
		body,
	)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orm-queries/users/transaction status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode create user response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("create user response success = false, want true")
	}

	var persisted entity.User
	query := app.Database().Table("users").Where("id = ?", expected.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read persisted user: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != expected.ID ||
		persisted.Email != expected.Email ||
		persisted.DepartmentID != expected.DepartmentID ||
		persisted.IsActive != expected.IsActive {
		t.Fatalf("persisted user = %+v, does not match create request", persisted)
	}
	if persisted.Password == "" || persisted.Password == expected.Password {
		t.Fatal("persisted password is empty or plaintext; want a hashed password")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("persisted createdAt is zero, want transaction to store a creation timestamp")
	}

	duplicate := expected
	duplicate.ID = "user-transaction-e2e-duplicate"
	duplicate.Password = "another-test-password"
	duplicateBody, err := json.Marshal(duplicate)
	if err != nil {
		t.Fatalf("encode duplicate user request: %v", err)
	}
	duplicateResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users/transaction",
		duplicateBody,
	)
	if duplicateResponse.StatusCode < http.StatusBadRequest {
		t.Fatalf("POST duplicate user status = %d, want 4xx or 5xx; body: %s",
			duplicateResponse.StatusCode, duplicateResponse.Body)
	}

	var persistedAfterDuplicate entity.User
	query = app.Database().Table("users").Where("id = ?", expected.ID).Find(&persistedAfterDuplicate)
	if query.Error != nil {
		t.Fatalf("read original user after duplicate create: %s: %v", query.Message, query.Error)
	}
	if persistedAfterDuplicate.Password != persisted.Password {
		t.Fatal("failed duplicate create changed the original user's password")
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orm-queries/users/"+expected.ID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orm-queries/users/%s status = %d, want %d; body: %s",
			expected.ID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}

	var getResult struct {
		Data entity.UserResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &getResult); err != nil {
		t.Fatalf("decode user response: %v", err)
	}
	if getResult.Data.ID != expected.ID ||
		getResult.Data.Email != expected.Email ||
		getResult.Data.DepartmentID != expected.DepartmentID ||
		getResult.Data.IsActive != expected.IsActive {
		t.Fatalf("user response = %+v, does not match create request", getResult.Data)
	}
}

func TestCreateUserWithSelectedColumnsE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	expected := entity.User{
		ID:           "user-selected-columns-e2e-01",
		Email:        "selected-columns-e2e@example.test",
		Password:     "plain-text-selected-columns-password",
		DepartmentID: "department-selected-columns-e2e-01",
		IsActive:     true,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("encode create user request: %v", err)
	}

	createResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users/selected-columns",
		body,
	)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orm-queries/users/selected-columns status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode create user response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("create user response success = false, want true")
	}

	var persisted entity.User
	query := app.Database().Table("users").Where("id = ?", expected.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read persisted user: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != expected.ID ||
		persisted.Email != expected.Email ||
		persisted.DepartmentID != expected.DepartmentID ||
		persisted.IsActive != expected.IsActive {
		t.Fatalf("persisted user = %+v, does not match create request", persisted)
	}
	if persisted.Password == "" || persisted.Password == expected.Password {
		t.Fatal("persisted password is empty or plaintext; want a hashed password")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("persisted createdAt is zero, want selected-columns create to store a timestamp")
	}
	if persisted.UpdatedAt != nil || persisted.DeletedAt != nil {
		t.Fatalf("unselected nullable timestamps should retain DB defaults; got updatedAt=%v deletedAt=%v",
			persisted.UpdatedAt, persisted.DeletedAt)
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orm-queries/users/"+expected.ID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orm-queries/users/%s status = %d, want %d; body: %s",
			expected.ID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}

	var getResult struct {
		Data entity.UserResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &getResult); err != nil {
		t.Fatalf("decode user response: %v", err)
	}
	if getResult.Data.ID != expected.ID ||
		getResult.Data.Email != expected.Email ||
		getResult.Data.DepartmentID != expected.DepartmentID ||
		getResult.Data.IsActive != expected.IsActive {
		t.Fatalf("user response = %+v, does not match create request", getResult.Data)
	}
}

func TestUpsertUserE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	initial := entity.User{
		ID:           "user-upsert-e2e-01",
		Email:        "upsert-e2e@example.test",
		Password:     "initial-upsert-password",
		DepartmentID: "department-upsert-before",
		IsActive:     true,
	}
	createUpsertUser(t, testServer, initial)

	updated := entity.User{
		ID:           "user-upsert-e2e-02",
		Email:        initial.Email,
		Password:     "updated-upsert-password",
		DepartmentID: "department-upsert-after",
		IsActive:     false,
	}
	createUpsertUser(t, testServer, updated)

	var persisted entity.User
	query := app.Database().Table("users").Where("email = ?", updated.Email).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read upserted user: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != initial.ID ||
		persisted.Email != updated.Email ||
		persisted.DepartmentID != updated.DepartmentID ||
		persisted.IsActive != updated.IsActive {
		t.Fatalf("persisted user after upsert = %+v, want ID %q retained, email %q, departmentID %q, isActive %t",
			persisted, initial.ID, updated.Email, updated.DepartmentID, updated.IsActive)
	}
	if persisted.Password == "" || persisted.Password == updated.Password {
		t.Fatal("upserted password is empty or plaintext; want a hashed password")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("upserted createdAt is zero, want a creation timestamp")
	}

	var users entity.UserArray
	query = app.Database().Table("users").Where("email = ?", updated.Email).Find(&users)
	if query.Error != nil {
		t.Fatalf("count records for upserted email: %s: %v", query.Message, query.Error)
	}
	if len(users) != 1 {
		t.Fatalf("records with upserted email = %d, want 1", len(users))
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orm-queries/users/"+initial.ID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orm-queries/users/%s status = %d, want %d; body: %s",
			initial.ID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}
	var getResult struct {
		Data entity.UserResponse `json:"data"`
	}
	if err := json.Unmarshal(getResponse.Body, &getResult); err != nil {
		t.Fatalf("decode upserted user response: %v", err)
	}
	if getResult.Data.ID != initial.ID ||
		getResult.Data.Email != updated.Email ||
		getResult.Data.DepartmentID != updated.DepartmentID ||
		getResult.Data.IsActive != updated.IsActive {
		t.Fatalf("user response after upsert = %+v, does not match upserted values", getResult.Data)
	}
}

func TestCreateUserReturningResultE2E(t *testing.T) {
	testServer, app := newUserTestServer(t)
	expected := entity.User{
		ID:           "user-returning-e2e-01",
		Email:        "returning-e2e@example.test",
		Password:     "plain-text-returning-password",
		DepartmentID: "department-returning-e2e-01",
		IsActive:     true,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("encode create user request: %v", err)
	}

	createResponse := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users/returning",
		body,
	)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orm-queries/users/returning status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool                `json:"success"`
		Data    entity.UserResponse `json:"data"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode create user response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("create user response success = false, want true")
	}
	if createResult.Data.ID != expected.ID {
		t.Fatalf("returned user ID = %q, want %q", createResult.Data.ID, expected.ID)
	}

	var persisted entity.User
	query := app.Database().Table("users").Where("id = ?", expected.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read persisted user: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != expected.ID ||
		persisted.Email != expected.Email ||
		persisted.DepartmentID != expected.DepartmentID ||
		persisted.IsActive != expected.IsActive {
		t.Fatalf("persisted user = %+v, does not match create request", persisted)
	}
	if persisted.Password == "" || persisted.Password == expected.Password {
		t.Fatal("persisted password is empty or plaintext; want a hashed password")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("persisted createdAt is zero, want a creation timestamp")
	}
	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(createResponse.Body, &responseFields); err != nil {
		t.Fatalf("decode response fields: %v", err)
	}
	var responseData map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseData); err != nil {
		t.Fatalf("decode returned user fields: %v", err)
	}
	if _, exposesPassword := responseData["password"]; exposesPassword {
		t.Fatal("create response exposes password; user response must not include it")
	}

	getResponse := sendRequest(
		t,
		testServer,
		http.MethodGet,
		"/examples/orm-queries/users/"+expected.ID,
		nil,
	)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET /examples/orm-queries/users/%s status = %d, want %d; body: %s",
			expected.ID, getResponse.StatusCode, http.StatusOK, getResponse.Body)
	}
}

func createUpsertUser(t *testing.T, testServer *server.Server, user entity.User) {
	t.Helper()

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("encode upsert user request: %v", err)
	}
	response := sendRequest(
		t,
		testServer,
		http.MethodPost,
		"/examples/orm-queries/users/upsert",
		body,
	)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/orm-queries/users/upsert status = %d, want 2xx; body: %s",
			response.StatusCode, response.Body)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode upsert response: %v", err)
	}
	if !result.Success {
		t.Fatal("upsert response success = false, want true")
	}
}

func newUserTestServer(t *testing.T) (*server.Server, *gorest.App) {
	t.Helper()

	dataDir := t.TempDir()
	copyUserMigration(t, dataDir)

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dataDir); err != nil {
		t.Fatalf("change to temporary working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	logConfig := gorestlog.ConfigDefault
	logConfig.LogDir = filepath.Join(dataDir, "logs")
	databaseConfig := database.Config{
		Driver:            database.SQLITE,
		Name:              "user-e2e",
		DatabaseDirectory: dataDir,
		Migration:         true,
		Enabled:           true,
	}
	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter-user-e2e",
		Version:     "test",
		Environment: "Test",
		Timezone:    "UTC",
	})
	app.Use(
		gorestlog.Run(logConfig),
		database.Run(databaseConfig),
	)
	t.Cleanup(func() {
		if message, _, err := app.Close(); err != nil {
			t.Errorf("close test application: %s: %v", message, err)
		}
	})
	if err := app.RunDatabaseStartupTasks(); err != nil {
		t.Fatalf("run SQLite user migration: %v", err)
	}

	testServer, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("create test HTTP server: %v", err)
	}
	userModel := model.NewUserModel(app)
	userLibrary := library.NewUserLibrary(userModel)
	userHandler := handler.NewUserHandler(app, userLibrary)
	router.RegisterUserRoutes(testServer.Group("/examples"), userHandler)

	return testServer, app
}

func copyUserMigration(t *testing.T, dataDir string) {
	t.Helper()

	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test source file")
	}
	repositoryRoot := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "../.."))
	source := filepath.Join(repositoryRoot, "app", "database", "migrations", "sqlite", "002_create_users_table.sql")
	destination := filepath.Join(dataDir, "migrations", "sqlite", filepath.Base(source))

	migration, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read user migration: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatalf("create temporary migration directory: %v", err)
	}
	if err := os.WriteFile(destination, migration, 0o600); err != nil {
		t.Fatalf("write temporary user migration: %v", err)
	}
}

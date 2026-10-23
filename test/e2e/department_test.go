package e2e_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/handler"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest-api-starter/app/router"
	"github.com/elsyahtech/gorest/database"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3"
)

func newDepartmentTestServer(t *testing.T) *server.Server {
	t.Helper()

	dataDir := t.TempDir()
	copyDepartmentMigration(t, dataDir)

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
		Name:              "department-e2e",
		DatabaseDirectory: dataDir,
		Migration:         true,
		Enabled:           true,
	}

	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter-e2e",
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
		t.Fatalf("run SQLite department migration: %v", err)
	}

	testServer, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("create test HTTP server: %v", err)
	}
	departmentModel := model.NewDepartmentModel(app)
	departmentLibrary := library.NewDepartmentLibrary(departmentModel)
	departmentHandler := handler.NewDepartmentHandler(app, departmentLibrary)
	router.RegisterDepartmentRoutes(testServer.Group("/examples"), departmentHandler)

	return testServer
}

type testHTTPResponse struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

func sendRequest(t *testing.T, testServer *server.Server, method, path string, body []byte) testHTTPResponse {
	t.Helper()

	request, err := http.NewRequestWithContext(
		context.Background(),
		method,
		path,
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("create %s request: %v", method, err)
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := testServer.SRV.Test(request, fiber.TestConfig{Timeout: 10 * time.Second})
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Error(err)
		}
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read %s response: %v", method, err)
	}

	return testHTTPResponse{
		StatusCode: response.StatusCode,
		Body:       responseBody,
		Header:     response.Header.Clone(),
	}
}

func copyDepartmentMigration(t *testing.T, dataDir string) {
	t.Helper()

	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test source file")
	}
	repositoryRoot := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "../.."))
	source := filepath.Join(repositoryRoot, "app/database/migrations/sqlite/001_create_departments_table.sql")
	destination := filepath.Join(dataDir, "migrations", "sqlite", filepath.Base(source))

	migration, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read department migration: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatalf("create temporary migration directory: %v", err)
	}
	if err := os.WriteFile(destination, migration, 0o600); err != nil {
		t.Fatalf("write temporary department migration: %v", err)
	}
}

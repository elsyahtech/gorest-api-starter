package e2e_test

import (
	"fmt"
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

func newOrderTestServer(t *testing.T) (*server.Server, string) {
	testServer, userIDs := newOrderTestServerWithUsers(t, 1)
	return testServer, userIDs[0]
}

func newOrderTestServerWithUsers(t *testing.T, userCount int) (*server.Server, []string) {
	t.Helper()

	dataDir := t.TempDir()
	copyOrderMigrations(t, dataDir)

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
		Name:              "order-e2e",
		DatabaseDirectory: dataDir,
		Migration:         true,
		Enabled:           true,
	}
	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter-order-e2e",
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
		t.Fatalf("run SQLite user and order migrations: %v", err)
	}

	userIDs := make([]string, 0, userCount)
	for i := 0; i < userCount; i++ {
		user := entity.User{
			ID:           fmt.Sprintf("user-order-e2e-%02d", i+1),
			Email:        fmt.Sprintf("order-e2e-%02d@example.test", i+1),
			Password:     "test-password",
			DepartmentID: "department-order-e2e-01",
			IsActive:     true,
		}
		if message, _, err := model.NewUserModel(app).CreateUser(&user); err != nil {
			t.Fatalf("create order test user %q: %s: %v", user.ID, message, err)
		}
		userIDs = append(userIDs, user.ID)
	}

	testServer, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("create test HTTP server: %v", err)
	}
	orderModel := model.NewOrderModel(app)
	orderLibrary := library.NewOrderLibrary(orderModel)
	orderHandler := handler.NewOrderHandler(app, orderLibrary)
	router.RegisterOrderRoutes(testServer.Group("/examples"), orderHandler)

	return testServer, userIDs
}

func copyOrderMigrations(t *testing.T, dataDir string) {
	t.Helper()

	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test source file")
	}
	repositoryRoot := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "../.."))
	migrationDir := filepath.Join(dataDir, "migrations", "sqlite")
	if err := os.MkdirAll(migrationDir, 0o755); err != nil {
		t.Fatalf("create temporary migration directory: %v", err)
	}

	for _, name := range []string{"002_create_users_table.sql", "003_create_orders_table.sql"} {
		source := filepath.Join(repositoryRoot, "app", "database", "migrations", "sqlite", name)
		migration, err := os.ReadFile(source)
		if err != nil {
			t.Fatalf("read migration %q: %v", name, err)
		}
		destination := filepath.Join(migrationDir, name)
		if err := os.WriteFile(destination, migration, 0o600); err != nil {
			t.Fatalf("write migration %q: %v", name, err)
		}
	}
}

package e2e_test

import (
	"fmt"
	"os"
	"path/filepath"
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
	databaseConfig := e2eDatabaseConfig(t, dataDir, "order-e2e")
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
		t.Fatalf("run %s user and order migrations: %v", databaseConfig.Driver, err)
	}
	if databaseConfig.Driver != database.SQLITE {
		resetSQLTestTables(t, app)
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
	for _, name := range []string{"002_create_users_table.sql", "003_create_orders_table.sql"} {
		copyE2EMigration(t, dataDir, name)
	}
}

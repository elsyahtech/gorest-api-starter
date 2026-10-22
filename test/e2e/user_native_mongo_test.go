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
	"github.com/elsyahtech/gorest/database"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/server"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestCreateUserWithNativeMongoQueryE2E(t *testing.T) {
	testServer, app, databaseName := newMongoUserTestServer(t)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, _, err := app.Database().RawQueryMongo(ctx, func(ctx context.Context, client *mongo.Client) error {
			return client.Database(databaseName).Drop(ctx)
		}); err != nil {
			t.Errorf("drop MongoDB E2E database: %v", err)
		}
	})

	expected := entity.User{
		ID:           "user-native-mongo-e2e-01",
		Email:        "native-mongo-e2e@example.test",
		Password:     "plain-text-native-mongo-password",
		DepartmentID: "department-native-mongo-e2e-01",
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
		"/examples/native-queries/users/mongodb",
		body,
	)
	if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("POST /examples/native-queries/users/mongodb status = %d, want 2xx; body: %s",
			createResponse.StatusCode, createResponse.Body)
	}

	var createResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(createResponse.Body, &createResult); err != nil {
		t.Fatalf("decode native MongoDB create response: %v", err)
	}
	if !createResult.Success {
		t.Fatal("native MongoDB create response success = false, want true")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var persistedUsers entity.UserArray
	if _, _, err := app.Database().RawQueryMongo(ctx, func(ctx context.Context, client *mongo.Client) error {
		cursor, err := client.Database(databaseName).
			Collection("users").
			Find(ctx, map[string]string{"id": expected.ID})
		if err != nil {
			return err
		}
		defer func() {
			if err := cursor.Close(ctx); err != nil {
				t.Error(err)
			}
		}()

		_, _, err = app.Database().DecodeMongoCursor(ctx, cursor, &persistedUsers)
		return err
	}); err != nil {
		t.Fatalf("read persisted MongoDB user: %v", err)
	}
	if len(persistedUsers) != 1 {
		t.Fatalf("persisted MongoDB user count = %d, want 1", len(persistedUsers))
	}
	persisted := persistedUsers[0]
	if persisted.ID != expected.ID ||
		persisted.Email != expected.Email ||
		persisted.DepartmentID != expected.DepartmentID ||
		persisted.IsActive != expected.IsActive {
		t.Fatalf("persisted MongoDB user = %+v, does not match request", persisted)
	}
	if persisted.Password == "" || persisted.Password == expected.Password {
		t.Fatal("persisted MongoDB password is empty or plaintext; want a hashed password")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("persisted MongoDB createdAt is zero, want a creation timestamp")
	}
}

func newMongoUserTestServer(t *testing.T) (*server.Server, *gorest.App, string) {
	t.Helper()
	if os.Getenv("GOREST_MONGO_E2E") != "1" {
		t.Skip("set GOREST_MONGO_E2E=1 to run the MongoDB-backed E2E test")
	}

	host := os.Getenv("GOREST_MONGO_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("GOREST_MONGO_PORT")
	if port == "" {
		port = "27018"
	}

	dataDir := t.TempDir()
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dataDir); err != nil {
		t.Fatalf("change to temporary working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDir); err != nil {
			t.Errorf("restore temporary working directory: %v", err)
		}
	})

	logConfig := gorestlog.ConfigDefault
	logConfig.LogDir = filepath.Join(dataDir, "logs")
	databaseName := fmt.Sprintf("gorest_user_e2e_%x", time.Now().UnixNano())
	databaseConfig := database.Config{
		Driver:            database.MONGO,
		Host:              host,
		Port:              port,
		Name:              databaseName,
		TLS:               "false",
		DatabaseDirectory: dataDir,
		Migration:         false,
		Seeder:            false,
		Enabled:           true,
		Timeout:           5 * time.Second,
	}
	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter-user-mongo-e2e",
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
			t.Errorf("close MongoDB test application: %s: %v", message, err)
		}
	})
	if err := app.RunDatabaseStartupTasks(); err != nil {
		t.Fatalf("initialize MongoDB test application: %v", err)
	}

	testServer, _, err := server.New(&server.ConfigDefault)
	if err != nil {
		t.Fatalf("create test HTTP server: %v", err)
	}
	userModel := model.NewUserModel(app)
	userLibrary := library.NewUserLibrary(userModel)
	userHandler := handler.NewUserHandler(app, userLibrary)
	router.RegisterUserRoutes(testServer.Group("/examples"), userHandler)

	return testServer, app, databaseName
}

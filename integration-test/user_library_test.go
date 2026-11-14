//go:build integration

package integration_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest-api-starter/app/library"
	"github.com/elsyahtech/gorest-api-starter/app/model"
	"github.com/elsyahtech/gorest/database"
	"github.com/elsyahtech/gorest/helper"
	"github.com/elsyahtech/gorest/log"
	"github.com/google/uuid"
)

func TestScyllaUserLibraryAgainstDatabase(t *testing.T) {
	testConfig, err := loadScyllaTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	prepareScyllaIntegrationSchema(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load ScyllaDB library-test configuration: %v", err)
	}
	dbConfig.Migration = false
	dbConfig.Seeder = false

	app := gorest.New(gorest.ConfigDefault)
	app.Use(log.Run(log.ConfigDefault), database.Run(*dbConfig))
	t.Cleanup(func() {
		if _, _, err := app.Close(); err != nil {
			t.Errorf("close ScyllaDB library-test application: %v", err)
		}
	})

	preflightConfig := *dbConfig
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	testDB, message, err := database.New(ctx, &preflightConfig, gorest.DefaultTimezone)
	cancel()
	if err != nil {
		t.Fatalf("connect to ScyllaDB library-test keyspace (%s): %v", message, err)
	}
	t.Cleanup(func() {
		closeScyllaDatabase(t, testDB, &preflightConfig)
	})

	userID := uuid.NewString()
	email := "library+" + strings.ReplaceAll(userID, "-", "") + "@example.test"
	cleanupScyllaUser(t, testDB, &preflightConfig, userID)
	t.Cleanup(func() {
		cleanupScyllaUser(t, testDB, &preflightConfig, userID)
	})

	userLibrary := library.NewUserLibrary(app, model.NewUserModel(app))
	plainPassword := "integration-library-password"
	user := &entity.User{
		ID:       userID,
		Email:    email,
		Password: plainPassword,
		IsActive: true,
	}
	message, status, err := userLibrary.CreateUser(user)
	if err != nil || status != http.StatusOK || message == "" {
		t.Fatalf("CreateUser() = (%q, %d, %v), want successful ScyllaDB persistence", message, status, err)
	}

	var storedPassword string
	queryContext, cancelQuery := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelQuery()
	if err := testDB.Scylla.Query("SELECT password FROM users WHERE id = ?", userID).
		WithContext(queryContext).Scan(&storedPassword); err != nil {
		t.Fatalf("read password persisted by UserLibrary: %v", err)
	}
	if storedPassword == plainPassword || storedPassword == "" {
		t.Fatal("UserLibrary persisted the plaintext password")
	}
	passwordMatches, _, _, err := helper.VerifyArgon(plainPassword, storedPassword)
	if err != nil || !passwordMatches {
		t.Fatalf("verify ScyllaDB password hash: matches=%t, err=%v", passwordMatches, err)
	}

	users, message, status, err := userLibrary.GetUsers()
	if err != nil || status != http.StatusOK || message == "" || users == nil {
		t.Fatalf("GetUsers() = (nil:%t, %q, %d, %v), want ScyllaDB users", users == nil, message, status, err)
	}
	found := false
	for _, listedUser := range *users {
		if listedUser.ID == userID && listedUser.Email == email {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("GetUsers() did not return created test user %q", userID)
	}

	message, status, err = userLibrary.Login(nil, &entity.User{
		Email:    email,
		Password: "incorrect-integration-password",
	})
	if err == nil || status != http.StatusUnauthorized || message != "" {
		t.Fatalf("Login() with incorrect password = (%q, %d, %v), want unauthorized", message, status, err)
	}
}

func cleanupScyllaUser(t *testing.T, db *database.Database, cfg *database.Config, id string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, _, err := db.ExecCQL(ctx, cfg, "DELETE FROM users WHERE id = ?", id); err != nil {
		t.Errorf("clean up ScyllaDB library-test user %s: %v", id, err)
	}
}

func closeScyllaDatabase(t *testing.T, db *database.Database, cfg *database.Config) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, _, err := db.Close(ctx, cfg); err != nil {
		t.Errorf("close ScyllaDB integration connection: %v", err)
	}
}

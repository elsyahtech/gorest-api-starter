//go:build integration

package integration_test

import (
	"context"
	"database/sql"
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

func TestPostgreSQLUserLibraryAgainstDatabase(t *testing.T) {
	testConfig, err := loadPostgreSQLTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	applyPostgreSQLTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load PostgreSQL library-test configuration: %v", err)
	}
	dbConfig.Migration = true
	dbConfig.Seeder = false
	dbConfig.MaxIdleConnections = 0

	app := gorest.New(gorest.ConfigDefault)
	app.Use(log.Run(log.ConfigDefault), database.Run(*dbConfig))
	t.Cleanup(func() {
		if _, _, err := app.Close(); err != nil {
			t.Errorf("close PostgreSQL library-test application: %v", err)
		}
	})

	preflightConfig := *dbConfig
	preflightConfig.Migration = false
	preflightConfig.Seeder = false
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	testDB, message, err := database.New(ctx, &preflightConfig, gorest.DefaultTimezone)
	cancel()
	if err != nil {
		t.Fatalf("connect to PostgreSQL library-test database (%s): %v", message, err)
	}
	t.Cleanup(func() {
		if err := testDB.SQL.Close(); err != nil {
			t.Errorf("close PostgreSQL library-test connection: %v", err)
		}
	})

	userID := uuid.NewString()
	email := "library+" + strings.ReplaceAll(userID, "-", "") + "@example.test"
	cleanupUserByID(t, testDB.SQL, userID)
	t.Cleanup(func() {
		cleanupUserByID(t, testDB.SQL, userID)
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
		t.Fatalf("CreateUser() = (%q, %d, %v), want successful PostgreSQL persistence", message, status, err)
	}

	var storedPassword string
	if err := testDB.SQL.QueryRow("SELECT password FROM users WHERE id = $1", userID).Scan(&storedPassword); err != nil {
		t.Fatalf("read password persisted by UserLibrary: %v", err)
	}
	if storedPassword == plainPassword || storedPassword == "" {
		t.Fatal("UserLibrary persisted the plaintext password")
	}
	passwordMatches, _, _, err := helper.VerifyArgon(plainPassword, storedPassword)
	if err != nil || !passwordMatches {
		t.Fatalf("verify PostgreSQL password hash: matches=%t, err=%v", passwordMatches, err)
	}

	users, message, status, err := userLibrary.GetUsers()
	if err != nil || status != http.StatusOK || message == "" || users == nil {
		t.Fatalf("GetUsers() = (nil:%t, %q, %d, %v), want PostgreSQL users", users == nil, message, status, err)
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

func cleanupUserByID(t *testing.T, db interface {
	Exec(query string, args ...any) (sql.Result, error)
}, id string) {
	t.Helper()
	if _, err := db.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		t.Errorf("clean up PostgreSQL library-test user %s: %v", id, err)
	}
}

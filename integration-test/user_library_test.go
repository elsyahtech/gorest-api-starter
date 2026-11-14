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
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestMongoUserLibraryAgainstDatabase(t *testing.T) {
	testConfig, err := loadMongoTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	prepareMongoIntegrationSchema(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load MongoDB library-test configuration: %v", err)
	}
	dbConfig.Migration = false
	dbConfig.Seeder = false

	app := gorest.New(gorest.ConfigDefault)
	app.Use(log.Run(log.ConfigDefault), database.Run(*dbConfig))
	t.Cleanup(func() {
		if _, _, err := app.Close(); err != nil {
			t.Errorf("close MongoDB library-test application: %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	testDB, message, err := database.New(ctx, dbConfig, gorest.DefaultTimezone)
	cancel()
	if err != nil {
		t.Fatalf("connect to MongoDB library-test database (%s): %v", message, err)
	}
	t.Cleanup(func() {
		cleanupMongoTestConnection(t, testDB.Mongo)
	})
	collection := testDB.Mongo.Database(testConfig.database).Collection("users")
	requireMongoSeedData(t, collection)

	userID := uuid.NewString()
	email := "library+" + strings.ReplaceAll(userID, "-", "") + "@example.test"
	cleanupMongoUser(t, collection, userID, email)
	t.Cleanup(func() {
		cleanupMongoUser(t, collection, userID, email)
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
		t.Fatalf("CreateUser() = (%q, %d, %v), want successful MongoDB persistence", message, status, err)
	}

	var stored struct {
		Password string `bson:"password"`
	}
	if err := collection.FindOne(context.Background(), bson.M{"id": userID}).Decode(&stored); err != nil {
		t.Fatalf("read password persisted by UserLibrary: %v", err)
	}
	if stored.Password == plainPassword || stored.Password == "" {
		t.Fatal("UserLibrary persisted the plaintext password")
	}
	passwordMatches, _, _, err := helper.VerifyArgon(plainPassword, stored.Password)
	if err != nil || !passwordMatches {
		t.Fatalf("verify MongoDB password hash: matches=%t, err=%v", passwordMatches, err)
	}

	users, message, status, err := userLibrary.GetUsers()
	if err != nil || status != http.StatusOK || message == "" || users == nil {
		t.Fatalf("GetUsers() = (nil:%t, %q, %d, %v), want MongoDB users", users == nil, message, status, err)
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

func requireMongoSeedData(t *testing.T, collection *mongo.Collection) {
	t.Helper()
	ids := []string{
		"3746e32f-1bd6-4023-a608-f535696f5e01",
		"3746e32f-1bd6-4023-a608-f535696f5e02",
		"3746e32f-1bd6-4023-a608-f535696f5e03",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, id := range ids {
		count, err := collection.CountDocuments(ctx, bson.M{"id": id})
		if err != nil {
			t.Fatalf("verify MongoDB sample user %s: %v", id, err)
		}
		if count != 1 {
			t.Fatalf("MongoDB seeder created %d documents for sample ID %s, want exactly one", count, id)
		}
	}
}

func cleanupMongoUser(t *testing.T, collection *mongo.Collection, id, email string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := collection.DeleteOne(ctx, bson.M{"id": id, "email": email}); err != nil {
		t.Errorf("clean up MongoDB integration user %s: %v", id, err)
	}
}

func cleanupMongoTestConnection(t *testing.T, client *mongo.Client) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		t.Errorf("disconnect MongoDB integration test client: %v", err)
	}
}

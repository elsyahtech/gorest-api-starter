package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestUpdateUserWithNativeMongoQueryE2E(t *testing.T) {
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

	createdAt := time.Date(2025, time.January, 15, 10, 0, 0, 0, time.UTC)
	initial := entity.User{
		ID:           "user-update-native-mongo-e2e-01",
		Email:        "update-native-mongo-old@example.test",
		Password:     "original-hashed-password",
		DepartmentID: "department-update-native-mongo-old",
		IsActive:     true,
		CreatedAt:    createdAt,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, _, err := app.Database().RawQueryMongo(ctx, func(ctx context.Context, client *mongo.Client) error {
		_, err := client.Database(databaseName).Collection("users").InsertOne(ctx, bson.M{
			"id":            initial.ID,
			"email":         initial.Email,
			"password":      initial.Password,
			"department_id": initial.DepartmentID,
			"is_active":     initial.IsActive,
			"created_at":    initial.CreatedAt,
		})
		return err
	}); err != nil {
		t.Fatalf("seed MongoDB user: %v", err)
	}

	updated := entity.User{
		ID:           initial.ID,
		Email:        "update-native-mongo-new@example.test",
		Password:     "plain-text-native-mongo-updated-password",
		DepartmentID: "department-update-native-mongo-new",
		IsActive:     false,
	}
	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("encode native MongoDB update request: %v", err)
	}

	path := "/examples/native-queries/users/mongodb/" + initial.ID
	response := sendRequest(t, testServer, http.MethodPut, path, body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("PUT %s status = %d, want 2xx; body: %s",
			path, response.StatusCode, response.Body)
	}
	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode native MongoDB update response: %v", err)
	}
	if !result.Success {
		t.Fatal("native MongoDB update response success = false, want true")
	}

	var persistedUsers entity.UserArray
	if _, _, err := app.Database().RawQueryMongo(ctx, func(ctx context.Context, client *mongo.Client) error {
		cursor, err := client.Database(databaseName).Collection("users").
			Find(ctx, bson.M{"id": initial.ID})
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
		t.Fatalf("read updated MongoDB user: %v", err)
	}
	if len(persistedUsers) != 1 {
		t.Fatalf("updated MongoDB user count = %d, want 1", len(persistedUsers))
	}
	persisted := persistedUsers[0]
	if persisted.ID != updated.ID ||
		persisted.Email != updated.Email ||
		persisted.DepartmentID != updated.DepartmentID ||
		persisted.IsActive != initial.IsActive {
		t.Fatalf("persisted MongoDB user after update = %+v", persisted)
	}
	if persisted.Password == "" || persisted.Password == updated.Password {
		t.Fatal("persisted MongoDB password is empty or plaintext; want a hashed password")
	}
	matchesPassword, message, _, err := helper.VerifyArgon(updated.Password, persisted.Password)
	if err != nil {
		t.Fatalf("verify persisted password hash: %s: %v", message, err)
	}
	if !matchesPassword {
		t.Fatal("persisted MongoDB password hash does not match submitted password")
	}
	if !persisted.CreatedAt.Equal(createdAt) {
		t.Fatalf("persisted MongoDB createdAt = %s, want unchanged value %s",
			persisted.CreatedAt.Format(time.RFC3339Nano), createdAt.Format(time.RFC3339Nano))
	}
	if persisted.UpdatedAt == nil || persisted.UpdatedAt.IsZero() {
		t.Fatal("persisted MongoDB updatedAt is nil or zero")
	}

	missingID := "user-update-native-mongo-e2e-missing"
	missingBody, err := json.Marshal(entity.User{
		ID:       missingID,
		Email:    "update-native-mongo-missing@example.test",
		Password: "plain-text-missing-user-password",
	})
	if err != nil {
		t.Fatalf("encode missing MongoDB user update request: %v", err)
	}
	missingPath := "/examples/native-queries/users/mongodb/" + missingID
	missingResponse := sendRequest(t, testServer, http.MethodPut, missingPath, missingBody)
	if missingResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("PUT %s status = %d, want %d; body: %s",
			missingPath, missingResponse.StatusCode, http.StatusNotFound, missingResponse.Body)
	}
}

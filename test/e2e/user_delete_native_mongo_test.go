package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestDeleteUserWithNativeMongoQueryE2E(t *testing.T) {
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

	userToDelete := entity.User{
		ID:           "user-delete-native-mongo-e2e-01",
		Email:        "delete-native-mongo-e2e@example.test",
		Password:     "plain-text-delete-native-mongo-password",
		DepartmentID: "department-delete-native-mongo-e2e-01",
		IsActive:     true,
	}
	userToKeep := entity.User{
		ID:           "user-delete-native-mongo-e2e-02",
		Email:        "keep-native-mongo-e2e@example.test",
		Password:     "plain-text-keep-native-mongo-password",
		DepartmentID: "department-delete-native-mongo-e2e-02",
		IsActive:     true,
	}
	for _, user := range []entity.User{userToDelete, userToKeep} {
		body, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("encode MongoDB user request: %v", err)
		}
		createResponse := sendRequest(
			t,
			testServer,
			http.MethodPost,
			"/examples/native-queries/users/mongodb",
			body,
		)
		if createResponse.StatusCode < http.StatusOK || createResponse.StatusCode >= http.StatusMultipleChoices {
			t.Fatalf("POST MongoDB user %s status = %d, want 2xx; body: %s",
				user.ID, createResponse.StatusCode, createResponse.Body)
		}
	}

	deleteResponse := sendRequest(
		t,
		testServer,
		http.MethodDelete,
		"/examples/native-queries/users/mongodb/"+userToDelete.ID,
		nil,
	)
	if deleteResponse.StatusCode < http.StatusOK || deleteResponse.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("DELETE /examples/native-queries/users/mongodb/%s status = %d, want 2xx; body: %s",
			userToDelete.ID, deleteResponse.StatusCode, deleteResponse.Body)
	}

	var deleteResult struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(deleteResponse.Body, &deleteResult); err != nil {
		t.Fatalf("decode native MongoDB delete response: %v", err)
	}
	if !deleteResult.Success {
		t.Fatal("native MongoDB delete response success = false, want true")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var deletedCount int64
	var persistedUsers entity.UserArray
	if _, _, err := app.Database().RawQueryMongo(ctx, func(ctx context.Context, client *mongo.Client) error {
		collection := client.Database(databaseName).Collection("users")
		var err error
		deletedCount, err = collection.CountDocuments(ctx, map[string]string{"id": userToDelete.ID})
		if err != nil {
			return err
		}

		cursor, err := collection.Find(ctx, map[string]string{"id": userToKeep.ID})
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
		t.Fatalf("verify MongoDB users after delete: %v", err)
	}
	if deletedCount != 0 {
		t.Fatalf("deleted MongoDB user count = %d, want 0", deletedCount)
	}
	if len(persistedUsers) != 1 || persistedUsers[0].ID != userToKeep.ID {
		t.Fatalf("unaffected MongoDB users = %+v, want only user %q", persistedUsers, userToKeep.ID)
	}
}

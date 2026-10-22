package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestGetUsersWithNativeMongoQueryE2E(t *testing.T) {
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

	path := "/examples/native-queries/users/mongodb"
	emptyResponse := sendRequest(t, testServer, http.MethodGet, path, nil)
	if emptyResponse.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, emptyResponse.StatusCode, http.StatusOK, emptyResponse.Body)
	}
	var emptyResult struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(emptyResponse.Body, &emptyResult); err != nil {
		t.Fatalf("decode empty native MongoDB users response: %v", err)
	}
	if emptyResult.Data == nil || len(emptyResult.Data) != 0 {
		t.Fatalf("empty native MongoDB users data = %#v, want an empty array", emptyResult.Data)
	}

	createdAt := time.Now().UTC()
	activeUsers := []entity.User{
		{
			ID:           "user-native-mongo-read-e2e-01",
			Email:        "native-mongo-read-01@example.test",
			Password:     "hashed-test-password-01",
			DepartmentID: "department-native-mongo-read-e2e-a",
			IsActive:     true,
			CreatedAt:    createdAt,
		},
		{
			ID:           "user-native-mongo-read-e2e-02",
			Email:        "native-mongo-read-02@example.test",
			Password:     "hashed-test-password-02",
			DepartmentID: "department-native-mongo-read-e2e-b",
			IsActive:     true,
			CreatedAt:    createdAt,
		},
	}
	inactiveUser := entity.User{
		ID:           "user-native-mongo-read-e2e-inactive",
		Email:        "native-mongo-read-inactive@example.test",
		Password:     "hashed-test-password-inactive",
		DepartmentID: "department-native-mongo-read-e2e-inactive",
		IsActive:     false,
		CreatedAt:    createdAt,
	}
	documents := make([]any, 0, len(activeUsers)+1)
	for _, user := range append(activeUsers, inactiveUser) {
		documents = append(documents, bson.M{
			"id":            user.ID,
			"email":         user.Email,
			"password":      user.Password,
			"department_id": user.DepartmentID,
			"is_active":     user.IsActive,
			"created_at":    user.CreatedAt,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, _, err := app.Database().RawQueryMongo(ctx, func(ctx context.Context, client *mongo.Client) error {
		_, err := client.Database(databaseName).Collection("users").InsertMany(ctx, documents)
		return err
	}); err != nil {
		t.Fatalf("seed native MongoDB users: %v", err)
	}

	response := sendRequest(t, testServer, http.MethodGet, path, nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d; body: %s",
			path, response.StatusCode, http.StatusOK, response.Body)
	}
	var result struct {
		Data entity.UserResponseArray `json:"data"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatalf("decode native MongoDB users response: %v", err)
	}
	if len(result.Data) != len(activeUsers) {
		t.Fatalf("native MongoDB active user count = %d, want %d; data=%+v",
			len(result.Data), len(activeUsers), result.Data)
	}

	usersByID := make(map[string]entity.UserResponse, len(result.Data))
	for _, user := range result.Data {
		usersByID[user.ID] = user
	}
	for _, expected := range activeUsers {
		user, found := usersByID[expected.ID]
		if !found {
			t.Errorf("native MongoDB response is missing active user %q", expected.ID)
			continue
		}
		if user.Email != expected.Email ||
			user.DepartmentID != expected.DepartmentID ||
			!user.IsActive ||
			user.CreatedAt == "" {
			t.Errorf("native MongoDB user response = %+v, does not match fixture %+v",
				user, expected)
		}
	}
	if _, found := usersByID[inactiveUser.ID]; found {
		t.Fatalf("inactive MongoDB user %q was returned", inactiveUser.ID)
	}

	var responseFields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &responseFields); err != nil {
		t.Fatalf("decode native MongoDB response fields: %v", err)
	}
	var responseUsers []map[string]json.RawMessage
	if err := json.Unmarshal(responseFields["data"], &responseUsers); err != nil {
		t.Fatalf("decode native MongoDB user fields: %v", err)
	}
	for i, user := range responseUsers {
		if _, hasPassword := user["password"]; hasPassword {
			t.Errorf("native MongoDB user at index %d unexpectedly exposes password", i)
		}
	}
}

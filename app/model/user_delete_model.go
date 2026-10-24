package model

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// DeleteUserWithTransaction removes a user in a transaction.
func (model *UserModel) DeleteUserWithTransaction(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	database := model.App.Database()

	transaction, message, httpCode, err := database.BeginTx()
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	defer func() {
		if message, _, err := transaction.Rollback(); err != nil {
			model.App.Log(gorest.Map{"error": err}).Error(message)
		}
	}()

	var user entity.User

	query := transaction.Table("users").Where("id = ?", userID).Delete(&user)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

// DeleteUser removes the user with the given ID.
func (model *UserModel) DeleteUser(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	database := model.App.Database()

	var user entity.User

	query := database.Table("users").Where("id = ?", userID).Delete(&user)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return query.Message, query.HTTPCode, nil
}

// DeleteUserWithNativeSQL removes a user using a native SQL statement.
func (model *UserModel) DeleteUserWithNativeSQL(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	database := model.App.Database()

	transaction, message, httpCode, err := database.BeginTx()
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	defer func() {
		if message, _, err := transaction.Rollback(); err != nil {
			model.App.Log(gorest.Map{"error": err}).Error(message)
		}
	}()

	driver := model.App.GetDBConfig().Driver
	query := `DELETE FROM users WHERE id = ` + convertPH(driver, 1)

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	_, message, httpCode, err = transaction.ExecSQL(ctx, query, userID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	if _, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return message, httpCode, nil
}

// DeleteUserWithNativeMongoQuery removes a user using the MongoDB driver.
func (model *UserModel) DeleteUserWithNativeMongoQuery(userID string) (string, int, error) {
	if userID == "" {
		return "User ID is required", http.StatusBadRequest,
			errors.New("delete user mongo: ID is empty")
	}

	database := model.App.Database()

	col, message, httpCode, err := database.MongoCollection("users")
	if err != nil {
		return message, httpCode, fmt.Errorf("delete user mongo: get collection: %w", err)
	}

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	filter := bson.M{"id": userID}

	result, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return "Failed to delete user", http.StatusInternalServerError,
			fmt.Errorf("delete user mongo: delete: %w", err)
	}

	if result.DeletedCount == 0 {
		return "User not found", http.StatusNotFound,
			errors.New("delete user mongo: no document matched the given ID")
	}

	return "User deleted successfully", http.StatusOK, nil
}

// DeleteUserWithNativeCQL removes a user using a native ScyllaDB CQL statement.
func (model *UserModel) DeleteUserWithNativeCQL(userID string) (string, int, error) {
	if userID == "" {
		return "User ID is required", http.StatusBadRequest,
			errors.New("delete user scylla: ID is empty")
	}

	database := model.App.Database()

	query := `DELETE FROM users WHERE id = ? IF EXISTS`

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	applied, message, httpCode, err := database.MapScanCAS(ctx, query, userID)
	if err != nil {
		msg := fmt.Sprintf("Failed to delete user: %s", message)
		return msg, httpCode, fmt.Errorf("delete user scylla: %w", err)
	}

	if !applied {
		return "User not found", http.StatusNotFound,
			errors.New("delete user scylla: document not found or not applied")
	}

	return "User deleted successfully", httpCode, nil
}

// DeleteUserInRedis removes a user from Redis.
func (model *UserModel) DeleteUserInRedis(userID string) (string, int, error) {
	if userID == "" {
		return "User ID is required", http.StatusBadRequest,
			errors.New("delete user redis: ID is empty")
	}

	redis := model.App.Redis()

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	redisClient, message, httpCode, err := redis.QueryRedis()
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	key := "user:" + userID
	deleted, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		return "Failed to delete user from Redis", http.StatusInternalServerError,
			fmt.Errorf("delete user redis: %w", err)
	}

	if deleted == 0 {
		return "User not found in Redis", http.StatusNotFound,
			errors.New("delete user redis: no key found for the given ID")
	}

	return "User deleted successfully from Redis", httpCode, nil
}

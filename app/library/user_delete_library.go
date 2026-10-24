package library

import (
	"errors"
	"fmt"
	"net/http"
)

// DeleteUserWithTransaction removes a user in a transaction.
func (lib *UserLibrary) DeleteUserWithTransaction(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	message, httpCode, err := lib.UserModel.DeleteUserWithTransaction(userID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// DeleteUser removes the user with the given ID.
func (lib *UserLibrary) DeleteUser(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	message, httpCode, err := lib.UserModel.DeleteUser(userID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// DeleteUserWithNativeSQL removes a user through native SQL.
func (lib *UserLibrary) DeleteUserWithNativeSQL(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	message, httpCode, err := lib.UserModel.DeleteUserWithNativeSQL(userID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// DeleteUserWithNativeMongoQuery removes a user through MongoDB.
func (lib *UserLibrary) DeleteUserWithNativeMongoQuery(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	message, httpCode, err := lib.UserModel.DeleteUserWithNativeMongoQuery(userID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// DeleteUserWithNativeCQL removes a user through ScyllaDB.
func (lib *UserLibrary) DeleteUserWithNativeCQL(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	message, httpCode, err := lib.UserModel.DeleteUserWithNativeCQL(userID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// DeleteUserInRedis removes a user from Redis.
func (lib *UserLibrary) DeleteUserInRedis(userID string) (string, int, error) {
	if userID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("delete user: user ID is empty")
	}

	message, httpCode, err := lib.UserModel.DeleteUserInRedis(userID)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

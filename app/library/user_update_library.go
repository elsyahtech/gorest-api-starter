package library

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

// UpdateUserWithTransaction hashes the password and updates a user in a transaction.
func (lib *UserLibrary) UpdateUserWithTransaction(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.UpdateUserWithTransaction(userID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// UpdateUserWithSelectedColumns updates selected user fields.
func (lib *UserLibrary) UpdateUserWithSelectedColumns(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	message, httpCode, err := lib.UserModel.UpdateUserWithSelectedColumns(userID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// UpdateUserReturningResult hashes the password and returns the updated user.
func (lib *UserLibrary) UpdateUserReturningResult(userID string, req *entity.User) (*entity.User, string, int, error) {
	if userID == "" {
		return nil, "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return nil, "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	result, message, httpCode, err := lib.UserModel.UpdateUserReturningResult(userID, req)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// UpdateUsersInBulk hashes passwords, sets update timestamps, and updates users.
func (lib *UserLibrary) UpdateUsersInBulk(req *entity.UserArray) (string, int, error) {
	if req == nil || len(*req) == 0 {
		return "Users must not be empty", http.StatusBadRequest,
			errors.New("bulk update users: request is empty")
	}

	timeNow := time.Now()

	for i := range *req {
		user := &(*req)[i]

		hashed, msg, code, err := helper.GenerateArgon(user.Password)
		if err != nil {
			return msg, code, fmt.Errorf("add bulk users: hash password for %s: %w", user.Email, err)
		}

		user.Password = hashed
		user.UpdatedAt = &timeNow
	}

	message, httpCode, err := lib.UserModel.UpdateUsersInBulk(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// UpdateUserWithNativeSQL hashes the password and updates a user through native SQL.
func (lib *UserLibrary) UpdateUserWithNativeSQL(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.UpdateUserWithNativeSQL(userID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// UpdateUserWithNativeMongoQuery hashes the password and updates a user through MongoDB.
func (lib *UserLibrary) UpdateUserWithNativeMongoQuery(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.UpdateUserWithNativeMongoQuery(userID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// UpdateUserWithNativeCQL hashes the password and updates a user through ScyllaDB.
func (lib *UserLibrary) UpdateUserWithNativeCQL(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.UpdateUserWithNativeCQL(userID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// UpdateUser hashes the password and updates a user.
func (lib *UserLibrary) UpdateUser(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.UpdateUser(userID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// PatchUserActiveStatus applies the supported partial user update.
func (lib *UserLibrary) PatchUserActiveStatus(userID string, req *entity.UserActivePatch) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("patch user: user ID is empty")
	}

	if req == nil || req.IsActive == nil {
		return "Provide isActive to patch the user", http.StatusBadRequest,
			errors.New("patch user: isActive is required")
	}

	message, httpCode, err := lib.UserModel.PatchUserActiveStatus(userID, req)
	if err != nil {
		return message, httpCode, fmt.Errorf("patch user: %w", err)
	}

	return message, httpCode, nil
}

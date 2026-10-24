package library

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

// CreateUserWithTransaction hashes the password and creates a user in a transaction.
func (lib *UserLibrary) CreateUserWithTransaction(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.CreateUserWithTransaction(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// CreateUserWithSelectedColumns hashes the password and creates selected user fields.
func (lib *UserLibrary) CreateUserWithSelectedColumns(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.CreateUserWithSelectedColumns(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// UpsertUser hashes the password and creates or updates a user.
func (lib *UserLibrary) UpsertUser(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.UpsertUser(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// CreateUserReturningResult hashes the password and returns the created user.
func (lib *UserLibrary) CreateUserReturningResult(req *entity.User) (*entity.User, string, int, error) {
	if req == nil {
		return nil, "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	result, message, httpCode, err := lib.UserModel.CreateUserReturningResult(req)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// CreateUsersInBulk validates users, hashes their passwords, and creates them.
func (lib *UserLibrary) CreateUsersInBulk(req *entity.UserArray) (string, int, error) {
	if req == nil || len(*req) == 0 {
		return "Users must not be empty", http.StatusBadRequest,
			errors.New("add bulk users: request is empty")
	}

	timeNow := time.Now()
	seen := make(map[string]struct{}, len(*req))

	// Validate emails, hash passwords, and set creation timestamps.
	for i := range *req {
		user := &(*req)[i]

		if _, dup := seen[user.Email]; dup {
			return "Duplicate email in request", http.StatusBadRequest,
				fmt.Errorf("add bulk users: duplicate email in request: %s", user.Email)
		}

		seen[user.Email] = struct{}{}

		hashed, msg, code, err := helper.GenerateArgon(user.Password)
		if err != nil {
			return msg, code, fmt.Errorf("add bulk users: hash password for %s: %w", user.Email, err)
		}

		user.Password = hashed
		user.CreatedAt = timeNow
	}

	message, httpCode, err := lib.UserModel.CreateUsersInBulk(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// CreateUserWithNativeSQL hashes the password and creates a user through native SQL.
func (lib *UserLibrary) CreateUserWithNativeSQL(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.CreateUserWithNativeSQL(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// CreateUserWithNativeMongoQuery hashes the password and creates a user through MongoDB.
func (lib *UserLibrary) CreateUserWithNativeMongoQuery(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.CreateUserWithNativeMongoQuery(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// CreateUserWithNativeCQL hashes the password and creates a user through ScyllaDB.
func (lib *UserLibrary) CreateUserWithNativeCQL(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.CreateUserWithNativeCQL(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// CreateUser hashes the password and creates a user.
func (lib *UserLibrary) CreateUser(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.CreateUser(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

// CreateUserToRedis stores a user in Redis.
func (lib *UserLibrary) CreateUserToRedis(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	passwordHashed, message, httpCode, err := helper.GenerateArgon(req.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	req.Password = passwordHashed

	message, httpCode, err = lib.UserModel.CreateUserToRedis(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	return message, httpCode, nil
}

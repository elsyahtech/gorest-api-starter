package library

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

func (lib *UserLibrary) Login(ctx *gorest.Context, req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	// Get user by email from the database.
	user, message, httpCode, err := lib.UserModel.GetUserByEmail(req.Email)
	if err != nil {
		return message, httpCode, fmt.Errorf("get user by email: %w", err)
	}

	// Validate the password using Argon2 hashing algorithm.
	passwordMatch, message, httpCode, err := helper.VerifyArgon(req.Password, user.Password)
	if err != nil {
		return message, httpCode, fmt.Errorf("verify password: %w", err)
	}

	if !passwordMatch {
		return "", http.StatusUnauthorized, errors.New("ensure you have entered the correct email and password")
	}

	// Generate a unique token ID for the user using uuid.
	// TokenID also can be used from user.ID.
	// Can be adjustable based on your application needs.
	tokenID := uuid.New().String() // Generate a unique token ID for the user.

	// Generate token authentication and store it in cookies.
	message, httpCode, err = lib.App.Filter().CreateToken(*ctx, tokenID)
	if err != nil {
		return message, httpCode, fmt.Errorf("create token: %w", err)
	}

	return message, httpCode, nil
}

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

func (lib *UserLibrary) GetUsers() (*[]entity.User, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsers()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

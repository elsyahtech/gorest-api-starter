package model

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// GetUserByEmail retrieves a user from the database based on the provided email address.
// It returns the user entity, a message, an HTTP status code, and an error if any occurred during the process.
func (model *UserModel) GetUserByEmail(email string) (*entity.User, string, int, error) {
	if email == "" {
		return nil, "Email param must not be empty", http.StatusBadRequest,
			errors.New("get user by email: email is empty")
	}

	database := model.App.Database()

	var user entity.User

	query := database.Table("users").
		Select("*").
		Where("users.email = ? AND users.is_active = ?", email, true).
		Find(&user)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &user, query.Message, query.HTTPCode, nil
}

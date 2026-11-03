package model

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func (model *UserModel) CreateUser(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: payload is empty")
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

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := transaction.Table("users").Create(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

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

func (model *UserModel) GetUsers() (*[]entity.User, string, int, error) {
	database := model.App.Database()

	var users []entity.User

	query := database.Table("users").Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, query.Message, query.HTTPCode, nil
}

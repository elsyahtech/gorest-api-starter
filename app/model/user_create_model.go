package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// CreateUserWithTransaction stores a user in a transaction.
func (model *UserModel) CreateUserWithTransaction(req *entity.User) (string, int, error) {
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

// CreateUserWithSelectedColumns stores selected user fields.
func (model *UserModel) CreateUserWithSelectedColumns(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	database := model.App.Database()

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := database.
		Table("users").
		Select("id, email, password, department_id, is_active, created_at").
		Create(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return query.Message, query.HTTPCode, nil
}

// UpsertUser creates or updates a user using the email conflict key.
func (model *UserModel) UpsertUser(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	database := model.App.Database()

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := database.
		Table("users").
		Upsert("email").
		Create(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return query.Message, query.HTTPCode, nil
}

// CreateUserReturningResult stores a user and returns the created record.
func (model *UserModel) CreateUserReturningResult(req *entity.User) (*entity.User, string, int, error) {
	if req == nil {
		return nil, "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	database := model.App.Database()

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := database.
		Table("users").
		Return("id").
		Create(req)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return req, query.Message, query.HTTPCode, nil
}

// CreateUsersInBulk stores multiple users in a transaction.
func (model *UserModel) CreateUsersInBulk(req *entity.UserArray) (string, int, error) {
	if req == nil || len(*req) == 0 {
		return "Users must not be empty", http.StatusBadRequest,
			errors.New("add bulk users: request is empty")
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

	query := transaction.Table("users").Create(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

// CreateUserWithNativeSQL stores a user using a native SQL statement.
func (model *UserModel) CreateUserWithNativeSQL(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
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

	query := `
		INSERT INTO users (id, email, password, department_id, is_active, created_at)
		VALUES (` + ph(driver, 6) + `)
	`

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	isActive := convertBooleanToInt(req.IsActive)

	req.CreatedAt = time.Now()

	_, message, httpCode, err = transaction.ExecSQL(ctx, query,
		req.ID,
		req.Email,
		req.Password,
		req.DepartmentID,
		isActive,
		req.CreatedAt,
	)

	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	if _, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return message, httpCode, nil
}

// CreateUserWithNativeMongoQuery stores a user using the MongoDB driver.
func (model *UserModel) CreateUserWithNativeMongoQuery(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user mongo: request is empty")
	}

	database := model.App.Database()

	col, message, httpCode, err := database.MongoCollection("users")
	if err != nil {
		return message, httpCode, fmt.Errorf("add user mongo: get collection: %w", err)
	}

	req.CreatedAt = time.Now()

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	doc, message, httpCode, err := database.EncodeMongoCollection(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("encode user: %w", err)
	}

	if _, err = col.InsertOne(ctx, doc); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "User already exists", http.StatusConflict,
				fmt.Errorf("add user mongo: insert: %w", err)
		}

		return "Failed to save user", http.StatusInternalServerError,
			fmt.Errorf("add user mongo: insert: %w", err)
	}

	return "User created successfully", http.StatusOK, nil
}

// CreateUserWithNativeCQL stores a user using a native ScyllaDB CQL statement.
func (model *UserModel) CreateUserWithNativeCQL(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	database := model.App.Database()

	req.CreatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, password, department_id, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS
	`

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	applied, message, httpCode, err := database.MapScanCAS(ctx, query,
		req.ID,
		req.Email,
		req.Password,
		req.DepartmentID,
		req.IsActive,
		req.CreatedAt,
	)
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	if !applied {
		return "Failed to create user", http.StatusConflict, errors.New("data violates a primary-key or unique constraint")
	}

	return message, httpCode, nil
}

// CreateUser stores a user without an explicit transaction.
func (model *UserModel) CreateUser(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}

	database := model.App.Database()

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := database.Table("users").Create(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return query.Message, query.HTTPCode, nil
}

// CreateUserToRedis stores a user in Redis.
func (model *UserModel) CreateUserToRedis(req *entity.User) (string, int, error) {
	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("add user: request is empty")
	}
	if req.ID == "" {
		return "User ID must not be empty", http.StatusBadRequest,
			errors.New("add user to Redis: user ID is empty")
	}

	redis := model.App.Redis()
	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	redisClient, message, httpCode, err := redis.QueryRedis()
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	user := *req
	user.CreatedAt = time.Now()

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return "Failed to process user data", http.StatusInternalServerError, fmt.Errorf("failed to marshal user: %w", err)
	}

	redisKey := fmt.Sprintf("user:%s", user.ID)

	created, err := redisClient.SetNX(ctx, redisKey, jsonBytes, 0).Result()
	if err != nil {
		return "Failed to save user in Redis", http.StatusInternalServerError, fmt.Errorf("failed to create user in Redis: %w", err)
	}
	if !created {
		return "User already exists in Redis", http.StatusConflict,
			errors.New("create user in Redis: key already exists")
	}

	return "User created successfully in Redis", httpCode, nil
}

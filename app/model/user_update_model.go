package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// UpdateUserWithTransaction updates a user in a transaction.
func (model *UserModel) UpdateUserWithTransaction(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
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

	query := transaction.Table("users").Where("id = ?", userID).Update(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

// UpdateUser updates a user without an explicit transaction.
func (model *UserModel) UpdateUser(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	database := model.App.Database()

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := database.Table("users").Where("id = ?", userID).Update(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return query.Message, query.HTTPCode, nil
}

// UpdateUserWithSelectedColumns updates selected user fields.
func (model *UserModel) UpdateUserWithSelectedColumns(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	database := model.App.Database()

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := database.
		Table("users").
		Select("id, email, updated_at").
		Where("id = ?", userID).
		Update(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return query.Message, query.HTTPCode, nil
}

// UpdateUserReturningResult updates a user and returns the updated record.
func (model *UserModel) UpdateUserReturningResult(userID string, req *entity.User) (*entity.User, string, int, error) {
	if userID == "" {
		return nil, "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return nil, "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	database := model.App.Database()

	timeNow := time.Now()

	req.CreatedAt = timeNow

	query := database.
		Table("users").
		Where("id = ?", userID).
		Return("id, email, created_at").
		Update(req)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return req, query.Message, query.HTTPCode, nil
}

// UpdateUsersInBulk updates multiple users in a transaction.
func (model *UserModel) UpdateUsersInBulk(req *entity.UserArray) (string, int, error) {
	if req == nil || len(*req) == 0 {
		return "Users must not be empty", http.StatusBadRequest,
			errors.New("bulk update users: request is empty")
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

	query := transaction.Table("users").Update(req)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	if message, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return query.Message, query.HTTPCode, nil
}

// UpdateUserWithNativeSQL updates a user using a native SQL statement.
func (model *UserModel) UpdateUserWithNativeSQL(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}

	database := model.App.Database()

	transaction, message, httpCode, err := database.BeginTx()
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	defer func() {
		if _, _, err := transaction.Rollback(); err != nil {
			model.App.Log(gorest.Map{"error": err}).Error("transaction rollback failed")
		}
	}()

	timeNow := time.Now()

	req.UpdatedAt = &timeNow

	driver := model.App.GetDBConfig().Driver
	activeValue := "true"

	if driver == "oracle" || driver == "sqlserver" {
		activeValue = "1"
	}

	query := `
		UPDATE users
		SET email = ` + convertPH(driver, 1) + `,
		password = ` + convertPH(driver, 2) + `,
		department_id = ` + convertPH(driver, 3) + `,
		updated_at = ` + convertPH(driver, 4) + `
		WHERE id = ` + convertPH(driver, 5) + ` AND is_active = ` + activeValue + `
	`

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	_, message, httpCode, err = transaction.ExecSQL(ctx, query,
		req.Email,
		req.Password,
		req.DepartmentID,
		req.UpdatedAt,
		userID,
	)

	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	if _, httpCode, err = transaction.Commit(); err != nil {
		return message, httpCode, fmt.Errorf("commit transaction failed: %w", err)
	}

	return message, httpCode, nil
}

// UpdateUserWithNativeMongoQuery updates a user using the MongoDB driver.
func (model *UserModel) UpdateUserWithNativeMongoQuery(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user mongo: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user mongo: payload is empty")
	}

	database := model.App.Database()

	col, message, httpCode, err := database.MongoCollection("users")
	if err != nil {
		return message, httpCode, fmt.Errorf("update user mongo: get collection: %w", err)
	}

	timeNow := time.Now()

	req.UpdatedAt = &timeNow

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	doc, message, httpCode, err := database.EncodeMongoCollection(req)
	if err != nil {
		return message, httpCode, fmt.Errorf("encode user: %w", err)
	}

	filter := bson.M{"id": userID}

	update := bson.M{
		"$set": bson.M{
			"email":         doc["email"],
			"password":      doc["password"],
			"department_id": doc["department_id"],
			"updated_at":    doc["updated_at"],
		},
	}

	result, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "User data conflict", http.StatusConflict,
				fmt.Errorf("update user mongo: update: %w", err)
		}

		return "Failed to update user", http.StatusInternalServerError,
			fmt.Errorf("update user mongo: update: %w", err)
	}

	if result.MatchedCount == 0 {
		return "User not found", http.StatusNotFound,
			errors.New("update user mongo: no document matched the given ID")
	}

	return "User updated successfully", http.StatusOK, nil
}

// UpdateUserWithNativeCQL updates a user using a native ScyllaDB CQL statement.
func (model *UserModel) UpdateUserWithNativeCQL(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user scylla: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user scylla: payload is empty")
	}

	database := model.App.Database()
	timeNow := time.Now()

	req.UpdatedAt = &timeNow

	query := `
		UPDATE users 
		SET email = ?, password = ?, department_id = ?, is_active = ?, updated_at = ?
		WHERE id = ?
		IF EXISTS
	`

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	applied, message, httpCode, err := database.MapScanCAS(ctx, query,
		req.Email,
		req.Password,
		req.DepartmentID,
		req.IsActive,
		req.UpdatedAt,
		userID,
	)
	if err != nil {
		msg := fmt.Sprintf("Failed to update user: %s", message)

		return msg, httpCode, fmt.Errorf("update user scylla: %w", err)
	}

	if !applied {
		return "User not found", http.StatusNotFound,
			errors.New("update user scylla: document not found or not applied")
	}

	return "User updated successfully", httpCode, nil
}

// PatchUserActiveStatus updates only the user's active status.
func (model *UserModel) PatchUserActiveStatus(userID string, req *entity.UserActivePatch) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("patch user: user ID is empty")
	}

	if req == nil || req.IsActive == nil {
		return "Provide isActive to patch the user", http.StatusBadRequest,
			errors.New("patch user: isActive is required")
	}

	database := model.App.Database()
	now := time.Now()
	updatedUser := entity.User{
		IsActive:  *req.IsActive,
		UpdatedAt: &now,
	}

	query := database.
		Table("users").
		Select("is_active, updated_at").
		Where("id = ?", userID).
		Update(&updatedUser)
	if query.Error != nil {
		return query.Message, query.HTTPCode, fmt.Errorf("patch user active state: %w", query.Error)
	}

	return "User status updated successfully", query.HTTPCode, nil
}

// UpdateUserInRedis updates a user's data in Redis.
func (model *UserModel) UpdateUserInRedis(userID string, req *entity.User) (string, int, error) {
	if userID == "" {
		return "User ID param must not be empty", http.StatusBadRequest,
			errors.New("update user: user ID is empty")
	}

	if req == nil {
		return "User payload must not be empty", http.StatusBadRequest,
			errors.New("update user: payload is empty")
	}
	if req.ID != "" && req.ID != userID {
		return "User ID in the request must match the URL ID", http.StatusBadRequest,
			errors.New("update user in Redis: request ID does not match URL ID")
	}

	redis := model.App.Redis()
	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	redisClient, message, httpCode, err := redis.QueryRedis()
	if err != nil {
		return message, httpCode, fmt.Errorf("%w", err)
	}

	redisKey := fmt.Sprintf("user:%s", userID)
	var updateErr error
	for attempt := 0; attempt < 3; attempt++ {
		updateErr = redisClient.Watch(ctx, func(tx *goredis.Tx) error {
			jsonBytes, err := tx.Get(ctx, redisKey).Bytes()
			if err != nil {
				return err
			}

			var existing entity.User
			if err := json.Unmarshal(jsonBytes, &existing); err != nil {
				return fmt.Errorf("decode existing Redis user: %w", err)
			}

			updated := *req
			updated.ID = userID
			updated.CreatedAt = existing.CreatedAt
			updatedAt := time.Now()
			updated.UpdatedAt = &updatedAt

			updatedJSON, err := json.Marshal(updated)
			if err != nil {
				return fmt.Errorf("marshal updated Redis user: %w", err)
			}

			_, err = tx.TxPipelined(ctx, func(pipe goredis.Pipeliner) error {
				pipe.Set(ctx, redisKey, updatedJSON, 0)
				return nil
			})
			return err
		}, redisKey)

		if !errors.Is(updateErr, goredis.TxFailedErr) {
			break
		}
	}

	if errors.Is(updateErr, goredis.Nil) {
		return "User not found in Redis", http.StatusNotFound,
			fmt.Errorf("update user in Redis: key does not exist: %w", goredis.Nil)
	}
	if errors.Is(updateErr, goredis.TxFailedErr) {
		return "User changed during update; please retry", http.StatusConflict,
			fmt.Errorf("update user in Redis: transaction retry limit reached: %w", updateErr)
	}
	if updateErr != nil {
		return "Failed to update user in Redis", http.StatusInternalServerError,
			fmt.Errorf("update user in Redis: %w", updateErr)
	}

	return "User updated successfully in Redis", httpCode, nil
}

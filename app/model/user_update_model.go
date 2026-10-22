package model

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
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

	// Oracle.
	// query := `
	// 	UPDATE users
	// 	SET email = :1, password = :2, department_id = :3, updated_at = :4
	// 	WHERE id = :5 AND is_active = 1
	// `

	// SQL Server.
	// query := `
	// 	UPDATE users
	// 	SET email = @p1, password = @p2, department_id = @p3, updated_at = @p4
	// 	WHERE id = @p5 AND is_active = 1
	// `

	// MySQL and SQLite.
	// query := `
	//         UPDATE users
	//         SET email = ?, password = ?, department_id = ?, updated_at = ? WHERE id = ? AND is_active = true
	// `

	// PostgreSQL.
	query := `
		UPDATE users
		SET email = $1, password = $2, department_id = $3, updated_at = $4
		WHERE id = $5 AND is_active = true
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

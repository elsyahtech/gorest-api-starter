package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GetUsers returns all users.
func (model *UserModel) GetUsers() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	query := database.Table("users").Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, query.Message, query.HTTPCode, nil
}

// GetUsersWithAllColumns returns users with all columns selected.
func (model *UserModel) GetUsersWithAllColumns() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	query := database.Table("users").
		Select("*").
		Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, query.Message, query.HTTPCode, nil
}

// GetUsersWithSelectedColumns returns users with selected columns.
func (model *UserModel) GetUsersWithSelectedColumns() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	query := database.Table("users").
		Select("id, email").
		Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, "", query.HTTPCode, nil
}

// GetUsersOrderedByID returns users ordered by ID in descending order.
func (model *UserModel) GetUsersOrderedByID() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	query := database.Table("users").
		Select("*").
		OrderBy("id", "DESC").
		Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, "", query.HTTPCode, nil
}

// GetUserCountsByDepartment returns user counts grouped by department.
func (model *UserModel) GetUserCountsByDepartment(groupBy string) (*entity.GroupByResponseArray, string, int, error) {
	database := model.App.Database()

	var summaries entity.GroupByResponseArray

	query := database.
		Table("users").
		Select("department_id, COUNT(*) AS user_count").
		GroupBy(groupBy).
		Find(&summaries)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &summaries, "", query.HTTPCode, nil
}

// GetUsersWithDistinctDepartmentIDs returns users with distinct department IDs.
// Distinct must not be combined with Select in GOREST.
func (model *UserModel) GetUsersWithDistinctDepartmentIDs() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	query := database.Table("users").
		Distinct("department_id").
		Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, "", query.HTTPCode, nil
}

// GetDepartmentsByUserCount returns department groups with more than two users.
func (model *UserModel) GetDepartmentsByUserCount() (*entity.GroupByResponseArray, string, int, error) {
	database := model.App.Database()

	var summaries entity.GroupByResponseArray

	query := database.
		Table("users").
		Select("department_id, COUNT(*) AS user_count").
		GroupBy("department_id").
		Having("COUNT(*) > ?", 2).
		Find(&summaries)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &summaries, "", query.HTTPCode, nil
}

// GetUsersWithDepartment returns users joined with their departments.
func (model *UserModel) GetUsersWithDepartment() (*entity.JoinResponseArray, string, int, error) {
	database := model.App.Database()

	var responseWithJoin entity.JoinResponseArray

	query := database.
		Table("users").
		Select("users.*, departments.name as department_name, departments.description as department_description").
		Join("left join departments on users.department_id = departments.id").
		Find(&responseWithJoin)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &responseWithJoin, "", query.HTTPCode, nil
}

// GetUsersWithPreloadedDepartment returns users with their departments preloaded.
// GOREST supports Preload for one-to-one, one-to-many, and many-to-many relations.
// Use Preload for defined relationships; do not combine it with Join.
func (model *UserModel) GetUsersWithPreloadedDepartment() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	query := database.
		Table("users").
		Select("*").
		Preload("Department").
		Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, "", query.HTTPCode, nil
}

// GetUsersWithNativeSQL returns users queried through native SQL.
func (model *UserModel) GetUsersWithNativeSQL() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var (
		users                entity.UserArray
		updatedAt, deletedAt sql.NullTime
	)

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	query := `
	    SELECT id, email, password, department_id, is_active, created_at, updated_at, deleted_at
	    FROM users where deleted_at IS NULL
	`

	rows, message, httpCode, err := database.QuerySQL(ctx, query)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			model.App.Log(gorest.Map{
				"error": err,
			}).Error("failed to close SQL rows")
		}
	}()

	for rows.Next() {
		var user entity.User

		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.DepartmentID,
			&user.IsActive,
			&user.CreatedAt,
			&updatedAt,
			&deletedAt,
		); err != nil {
			message = "check if the column structure in the 'users' table matches the field sequence of the User DB Entity."

			return nil, message, http.StatusInternalServerError, fmt.Errorf("scan user: %w", err)
		}

		if updatedAt.Valid {
			user.UpdatedAt = &updatedAt.Time
		}

		if deletedAt.Valid {
			user.DeletedAt = &deletedAt.Time
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		message = "an error occurred while iterating through the result set. " +
			"Check network stability or database cursor/timeout configuration issues in ./app/config/database."

		return nil, message, http.StatusInternalServerError, fmt.Errorf("failed to get data into the entity struct: %w", err)
	}

	return &users, message, http.StatusOK, nil
}

// GetUsersWithNativeMongoQuery returns active users matching the MongoDB query.
func (model *UserModel) GetUsersWithNativeMongoQuery() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	coll, message, httpCode, err := database.MongoCollection("users")
	if err != nil {
		return nil, message, httpCode,
			fmt.Errorf("get users collection: %w", err)
	}

	cursor, err := coll.Find(ctx, bson.M{
		"is_active": true,
	})
	if err != nil {
		message = "failed to execute the MongoDB find operation; check the underlying error for details."

		return nil, message, http.StatusInternalServerError,
			fmt.Errorf("find MongoDB users: %w", err)
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			model.App.Log(gorest.Map{
				"error": err,
			}).Error("failed to close MongoDB cursor")
		}
	}()

	message, httpCode, err = database.DecodeMongoCursor(ctx, cursor, &users)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("decode MongoDB users: %w", err)
	}

	return &users, message, http.StatusOK, nil
}

// GetUsersWithNativeCQL returns users queried through native ScyllaDB CQL.
func (model *UserModel) GetUsersWithNativeCQL() (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	query := `
        SELECT id, email, password, department_id, is_active, created_at, updated_at, deleted_at
        FROM users
    `

	iter, message, httpCode, err := database.QueryCQL(ctx, query)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	var (
		id           string
		email        string
		password     string
		departmentID string
		isActive     bool
		createdAt    time.Time
		updatedAt    *time.Time
		deletedAt    *time.Time
	)

	for iter.Scan(
		&id,
		&email,
		&password,
		&departmentID,
		&isActive,
		&createdAt,
		&updatedAt,
		&deletedAt,
	) {
		var user entity.User

		user.ID = id
		user.Email = email
		user.Password = password
		user.DepartmentID = departmentID
		user.IsActive = isActive
		user.CreatedAt = createdAt
		user.UpdatedAt = updatedAt
		user.DeletedAt = deletedAt

		users = append(users, user)
	}

	if err := iter.Close(); err != nil {
		message := "an error occurred while iterating through the result set from ScyllaDB."

		model.App.Log(gorest.Map{"error": err}).Error(message)

		return nil, message, http.StatusInternalServerError, fmt.Errorf("failed to get data into entity struct: %w", err)
	}

	return &users, message, http.StatusOK, nil
}

// GetUsersByLimit returns at most limit users.
func (model *UserModel) GetUsersByLimit(limit int) (*entity.UserArray, string, int, error) {
	if limit == 0 {
		users := make(entity.UserArray, 0)
		return &users, "", http.StatusOK, nil
	}

	database := model.App.Database()

	var users entity.UserArray

	query := database.Table("users").
		Select("*").
		Limit(limit).
		Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, "", query.HTTPCode, nil
}

// GetUsersByLimitOffset returns users for the requested page.
func (model *UserModel) GetUsersByLimitOffset(limit, offset int) (*entity.UserArray, string, int, error) {
	database := model.App.Database()

	var users entity.UserArray

	query := database.Table("users").
		Select("*").
		OrderBy("users.id", "ASC").
		Limit(limit).
		Offset(offset).
		Find(&users)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &users, "", query.HTTPCode, nil
}

// GetUserByID returns the user with the given ID.
func (model *UserModel) GetUserByID(userID string) (*entity.User, string, int, error) {
	if userID == "" {
		return nil, "User ID param must not be empty", http.StatusBadRequest,
			errors.New("get user by ID: user ID is empty")
	}

	database := model.App.Database()

	var user entity.User

	query := database.Table("users").
		Select("*").
		Where("users.id = ?", userID).
		Find(&user)
	if query.Error != nil {
		return nil, query.Message, query.HTTPCode, fmt.Errorf("%w", query.Error)
	}

	return &user, "", query.HTTPCode, nil
}

// GetUsersFromRedis retrieves users stored in Redis.
func (model *UserModel) GetUsersFromRedis() (*entity.UserArray, string, int, error) {
	redis := model.App.Redis()
	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	redisClient, message, httpCode, err := redis.QueryRedis()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	var users entity.UserArray

	iter := redisClient.Scan(ctx, 0, "user:*", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		var user entity.User

		jsonData, err := redisClient.Get(ctx, key).Bytes()
		if err != nil {
			return nil, "Failed to retrieve user data from Redis", http.StatusInternalServerError,
				fmt.Errorf("failed to get user data from Redis for key %s: %w", key, err)
		}

		err = json.Unmarshal(jsonData, &user)
		if err != nil {
			return nil, "Failed to unmarshal user data from Redis", http.StatusInternalServerError,
				fmt.Errorf("failed to unmarshal user data from Redis for key %s: %w", key, err)
		}

		users = append(users, user)
	}

	if err := iter.Err(); err != nil {
		return nil, "Failed to scan user keys from Redis", http.StatusInternalServerError,
			fmt.Errorf("redis scan iterator error: %w", err)
	}

	return &users, "Users retrieved successfully from Redis", httpCode, nil
}

// GetUserFromRedisByID retrieves a user by ID from Redis.
func (model *UserModel) GetUserFromRedisByID(userID string) (*entity.User, string, int, error) {
	if userID == "" {
		return nil, "User ID param must not be empty", http.StatusBadRequest,
			errors.New("get user from Redis by ID: user ID is empty")
	}

	getRedis := model.App.Redis()
	ctx, cancel := model.App.NewContext(5 * time.Second)
	defer cancel()

	redisClient, message, httpCode, err := getRedis.QueryRedis()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	key := "user:" + userID
	jsonData, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, "User not found in Redis", http.StatusNotFound,
				fmt.Errorf("user with ID %s not found in Redis: %w", userID, err)
		}

		return nil, "Failed to retrieve user data from Redis", http.StatusInternalServerError,
			fmt.Errorf("failed to get user data from Redis for key %s: %w", key, err)
	}

	var user entity.User

	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		return nil, "Failed to unmarshal user data from Redis", http.StatusInternalServerError,
			fmt.Errorf("failed to unmarshal user data from Redis for key %s: %w", key, err)
	}

	return &user, "User retrieved successfully from Redis", httpCode, nil
}

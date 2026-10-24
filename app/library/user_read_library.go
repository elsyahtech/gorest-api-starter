package library

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

// GetUsers returns all users.
func (lib *UserLibrary) GetUsers() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsers()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersWithAllColumns returns users with all model columns selected.
func (lib *UserLibrary) GetUsersWithAllColumns() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersWithAllColumns()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersWithSelectedColumns returns users with selected model columns.
func (lib *UserLibrary) GetUsersWithSelectedColumns() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersWithSelectedColumns()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersOrderedByID returns users ordered by ID.
func (lib *UserLibrary) GetUsersOrderedByID() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersOrderedByID()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUserCountsByDepartment returns user counts grouped by department.
func (lib *UserLibrary) GetUserCountsByDepartment() (*entity.GroupByResponseArray, string, int, error) {
	groupBy := "department_id"

	result, message, httpCode, err := lib.UserModel.GetUserCountsByDepartment(groupBy)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersWithDistinctDepartmentIDs returns users with distinct department IDs.
func (lib *UserLibrary) GetUsersWithDistinctDepartmentIDs() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersWithDistinctDepartmentIDs()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetDepartmentsByUserCount returns department groups filtered by user count.
func (lib *UserLibrary) GetDepartmentsByUserCount() (*entity.GroupByResponseArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetDepartmentsByUserCount()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersWithDepartment returns users joined with their departments.
func (lib *UserLibrary) GetUsersWithDepartment() (*entity.JoinResponseArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersWithDepartment()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersWithPreloadedDepartment returns users with their departments preloaded.
func (lib *UserLibrary) GetUsersWithPreloadedDepartment() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersWithPreloadedDepartment()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersWithNativeSQL returns users queried through native SQL.
func (lib *UserLibrary) GetUsersWithNativeSQL() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersWithNativeSQL()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersWithNativeMongoQuery returns users matching the native MongoDB query.
func (lib *UserLibrary) GetUsersWithNativeMongoQuery() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersWithNativeMongoQuery()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersWithNativeCQL returns users queried through native ScyllaDB CQL.
func (lib *UserLibrary) GetUsersWithNativeCQL() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersWithNativeCQL()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersByLimit returns at most limit users.
func (lib *UserLibrary) GetUsersByLimit(limit int) (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersByLimit(limit)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersByLimitOffset returns users for the given page.
func (lib *UserLibrary) GetUsersByLimitOffset(limit, offset int) (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersByLimitOffset(limit, offset)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUserByID returns the user with the given ID.
func (lib *UserLibrary) GetUserByID(userID string) (*entity.User, string, int, error) {
	if userID == "" {
		return nil, "User ID param must not be empty", http.StatusBadRequest,
			errors.New("get user by ID: user ID is empty")
	}

	result, message, httpCode, err := lib.UserModel.GetUserByID(userID)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUsersFromRedis retrieves users stored in Redis.
func (lib *UserLibrary) GetUsersFromRedis() (*entity.UserArray, string, int, error) {
	result, message, httpCode, err := lib.UserModel.GetUsersFromRedis()
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

// GetUserFromRedisByID retrieves a user by ID from Redis.
func (lib *UserLibrary) GetUserFromRedisByID(userID string) (*entity.User, string, int, error) {
	if userID == "" {
		return nil, "User ID param must not be empty", http.StatusBadRequest,
			errors.New("get user from Redis by ID: user ID is empty")
	}

	result, message, httpCode, err := lib.UserModel.GetUserFromRedisByID(userID)
	if err != nil {
		return nil, message, httpCode, fmt.Errorf("%w", err)
	}

	return result, message, httpCode, nil
}

package entity

import "time"

// User represents a user entity.
type User struct {
	CreatedAt    time.Time   `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt    *time.Time  `gorest:"updated_at" json:"updatedAt"`
	DeletedAt    *time.Time  `gorest:"deleted_at" json:"deletedAt"`
	Department   *Department `gorest:"-" references:"ID" json:"department,omitempty"`
	ID           string      `gorest:"id, primary_key, required" json:"id"`
	Email        string      `gorest:"email, required, unique" json:"email"`
	Password     string      `gorest:"password, required" json:"password"`
	DepartmentID string      `gorest:"department_id, required" json:"departmentId"`
	IsActive     bool        `gorest:"is_active" json:"isActive"`
}

// UserActivePatch is a partial request DTO. A pointer lets the API distinguish
// an omitted isActive property from an explicit false value.
type UserActivePatch struct {
	IsActive *bool `json:"isActive"`
}

// UserArray is a collection of user records.
type UserArray []User

// UserResponse is the JSON representation of a user.
type UserResponse struct {
	CreatedAt    string      `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt    *string     `gorest:"updated_at" json:"updatedAt"`
	DeletedAt    *string     `gorest:"deleted_at" json:"deletedAt"`
	Department   *Department `gorest:"-" references:"ID" json:"department,omitempty"`
	ID           string      `gorest:"id, primary_key" json:"id"`
	Email        string      `gorest:"email" json:"email"`
	DepartmentID string      `gorest:"department_id" json:"departmentId"`
	IsActive     bool        `gorest:"is_active" json:"isActive"`
}

// UserResponseArray is a collection of user response DTOs.
type UserResponseArray []UserResponse

// GroupByResponse contains the user count for one department.
type GroupByResponse struct {
	DepartmentID string `gorest:"department_id" json:"departmentId"`
	UserCount    int64  `gorest:"user_count" json:"userCount"`
}

// GroupByResponseArray is a collection of department user counts.
type GroupByResponseArray []GroupByResponse

// JoinResponse contains a user and the selected fields from its department.
type JoinResponse struct {
	CreatedAt             string  `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt             *string `gorest:"updated_at" json:"updatedAt"`
	DeletedAt             *string `gorest:"deleted_at" json:"deletedAt"`
	ID                    string  `gorest:"id, primary_key" json:"id"`
	Email                 string  `gorest:"email" json:"email"`
	DepartmentID          string  `gorest:"department_id" json:"departmentId"`
	DepartmentName        string  `gorest:"department_name" json:"departmentName"`
	DepartmentDescription string  `gorest:"department_description" json:"departmentDescription"`
	IsActive              bool    `gorest:"is_active" json:"isActive"`
}

// JoinResponseArray is a collection of joined user and department records.
type JoinResponseArray []JoinResponse

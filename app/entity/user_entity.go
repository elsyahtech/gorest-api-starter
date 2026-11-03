package entity

import "time"

// User represents a user entity.
type User struct {
	CreatedAt time.Time  `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt *time.Time `gorest:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorest:"deleted_at" json:"deletedAt"`
	ID        string     `gorest:"id, primary_key, required" json:"id"`
	Email     string     `gorest:"email, required, unique" json:"email"`
	Password  string     `gorest:"password, required" json:"password"`
	IsActive  bool       `gorest:"is_active" json:"isActive"`
}

// UserResponse is the JSON representation of a user.
type UserResponse struct {
	CreatedAt string  `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt *string `gorest:"updated_at" json:"updatedAt"`
	DeletedAt *string `gorest:"deleted_at" json:"deletedAt"`
	ID        string  `gorest:"id, primary_key" json:"id"`
	Email     string  `gorest:"email" json:"email"`
	IsActive  bool    `gorest:"is_active" json:"isActive"`
}

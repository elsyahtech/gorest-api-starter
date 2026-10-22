package entity

import "time"

// Department represents a department database record.
type Department struct {
	CreatedAt   time.Time  `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt   *time.Time `gorest:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `gorest:"deleted_at" json:"deletedAt"`
	ID          string     `gorest:"id, primary_key, required" json:"id"`
	Name        string     `gorest:"name, required" json:"name"`
	Description string     `gorest:"description, required" json:"description"`
	IsActive    bool       `gorest:"is_active" json:"isActive"`
}

// DepartmentPatch is a request DTO for updating selected department fields.
type DepartmentPatch struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"isActive"`
}

// DepartmentArray is a collection of department records.
type DepartmentArray []Department

// DepartmentResponse is the JSON representation of a department.
type DepartmentResponse struct {
	CreatedAt   string  `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt   *string `gorest:"updated_at" json:"updatedAt"`
	DeletedAt   *string `gorest:"deleted_at" json:"deletedAt"`
	ID          string  `gorest:"id, primary_key" json:"id"`
	Name        string  `gorest:"name" json:"name"`
	Description string  `gorest:"description" json:"description"`
	IsActive    bool    `gorest:"is_active" json:"isActive"`
}

// DepartmentResponseArray is a collection of department response DTOs.
type DepartmentResponseArray []DepartmentResponse

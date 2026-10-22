package entity

import "time"

// Order represents an order database record.
type Order struct {
	CreatedAt   time.Time  `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt   *time.Time `gorest:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `gorest:"deleted_at" json:"deletedAt"`
	ID          string     `gorest:"id, primary_key, required" json:"id"`
	UserID      string     `gorest:"user_id, required" json:"userId"`
	OrderNumber string     `gorest:"order_number, required" json:"orderNumber"`
	Status      string     `gorest:"status, required" json:"status"`
	TotalAmount float64    `gorest:"total_amount, required" json:"totalAmount"`
	IsActive    bool       `gorest:"is_active" json:"isActive"`
}

// OrderPatch is a request DTO for updating selected order fields.
type OrderPatch struct {
	OrderNumber *string  `json:"orderNumber"`
	Status      *string  `json:"status"`
	TotalAmount *float64 `json:"totalAmount"`
	IsActive    *bool    `json:"isActive"`
}

// OrderArray is a collection of order records.
type OrderArray []Order

// OrderResponse is the JSON representation of an order.
type OrderResponse struct {
	CreatedAt   string  `gorest:"created_at, immutable" json:"createdAt"`
	UpdatedAt   *string `gorest:"updated_at" json:"updatedAt"`
	DeletedAt   *string `gorest:"deleted_at" json:"deletedAt"`
	ID          string  `gorest:"id, primary_key" json:"id"`
	UserID      string  `gorest:"user_id" json:"userId"`
	OrderNumber string  `gorest:"order_number" json:"orderNumber"`
	Status      string  `gorest:"status" json:"status"`
	TotalAmount float64 `gorest:"total_amount" json:"totalAmount"`
	IsActive    bool    `gorest:"is_active" json:"isActive"`
}

// OrderResponseArray is a collection of order response DTOs.
type OrderResponseArray []OrderResponse

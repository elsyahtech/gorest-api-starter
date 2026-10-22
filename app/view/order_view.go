package view

import (
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

// BuildOrderResponse converts an order into its JSON response DTO.
func BuildOrderResponse(response *entity.Order) *entity.OrderResponse {
	if response == nil {
		return nil
	}

	return &entity.OrderResponse{
		ID:          response.ID,
		UserID:      response.UserID,
		OrderNumber: response.OrderNumber,
		Status:      response.Status,
		TotalAmount: response.TotalAmount,
		IsActive:    response.IsActive,
		CreatedAt:   helper.ConvertTimeToString(response.CreatedAt),
		UpdatedAt:   helper.ConvertTimePointerToStringPointer(response.UpdatedAt),
		DeletedAt:   helper.ConvertTimePointerToStringPointer(response.DeletedAt),
	}
}

// BuildOrderResponseArray converts orders into response DTOs.
// It returns an empty array when responses is nil or empty.
func BuildOrderResponseArray(responses *entity.OrderArray) *entity.OrderResponseArray {
	if responses == nil || len(*responses) == 0 {
		return &entity.OrderResponseArray{}
	}

	res := make(entity.OrderResponseArray, 0, len(*responses))

	for _, item := range *responses {
		res = append(res, *BuildOrderResponse(&item))
	}

	return &res
}

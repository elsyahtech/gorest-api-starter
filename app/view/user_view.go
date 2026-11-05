package view

import (
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

// BuildUserResponse converts a user into its JSON response DTO.
func BuildUserResponse(response *entity.User) *entity.UserResponse {
	if response == nil {
		return nil
	}

	return &entity.UserResponse{
		ID:        response.ID,
		Email:     response.Email,
		IsActive:  response.IsActive,
		CreatedAt: helper.ConvertTimeToString(response.CreatedAt),
		UpdatedAt: helper.ConvertTimePointerToStringPointer(response.UpdatedAt),
		DeletedAt: helper.ConvertTimePointerToStringPointer(response.DeletedAt),
	}
}

// BuildUserResponseArray converts users into response DTOs.
// It returns an empty array when responses is nil or empty.
func BuildUserResponseArray(responses *[]entity.User) *[]entity.UserResponse {
	if responses == nil || len(*responses) == 0 {
		users := make([]entity.UserResponse, 0)

		return &users
	}

	res := make([]entity.UserResponse, 0, len(*responses))

	for _, item := range *responses {
		res = append(res, *BuildUserResponse(&item))
	}

	return &res
}

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
		ID:           response.ID,
		Email:        response.Email,
		DepartmentID: response.DepartmentID,
		Department:   response.Department,
		IsActive:     response.IsActive,
		CreatedAt:    helper.ConvertTimeToString(response.CreatedAt),
		UpdatedAt:    helper.ConvertTimePointerToStringPointer(response.UpdatedAt),
		DeletedAt:    helper.ConvertTimePointerToStringPointer(response.DeletedAt),
	}
}

// BuildUserResponseArray converts users into response DTOs.
// It returns an empty array when responses is nil or empty.
func BuildUserResponseArray(responses *entity.UserArray) *entity.UserResponseArray {
	if responses == nil || len(*responses) == 0 {
		return &entity.UserResponseArray{}
	}

	res := make(entity.UserResponseArray, 0, len(*responses))

	for _, item := range *responses {
		res = append(res, *BuildUserResponse(&item))
	}

	return &res
}

// BuildGroupByResponse returns the group-by response unchanged.
func BuildGroupByResponse(response *entity.GroupByResponse) *entity.GroupByResponse {
	if response == nil {
		return nil
	}

	return response
}

// BuildGroupByResponseArray copies group-by results into a response array.
// It returns an empty array when responses is nil or empty.
func BuildGroupByResponseArray(responses *entity.GroupByResponseArray) *entity.GroupByResponseArray {
	if responses == nil || len(*responses) == 0 {
		return &entity.GroupByResponseArray{}
	}

	res := make(entity.GroupByResponseArray, 0, len(*responses))

	for _, item := range *responses {
		res = append(res, *BuildGroupByResponse(&item))
	}

	return &res
}

// ToJoinResponse copies a joined user and department response.
func ToJoinResponse(response *entity.JoinResponse) *entity.JoinResponse {
	if response == nil {
		return nil
	}

	return &entity.JoinResponse{
		ID:                    response.ID,
		Email:                 response.Email,
		DepartmentID:          response.DepartmentID,
		DepartmentName:        response.DepartmentName,
		DepartmentDescription: response.DepartmentDescription,
		IsActive:              response.IsActive,
		CreatedAt:             response.CreatedAt,
		UpdatedAt:             response.UpdatedAt,
		DeletedAt:             response.DeletedAt,
	}
}

// ToJoinResponses copies joined user and department responses into an array.
// It returns an empty array when responses is nil or empty.
func ToJoinResponses(responses *entity.JoinResponseArray) *entity.JoinResponseArray {
	if responses == nil || len(*responses) == 0 {
		return &entity.JoinResponseArray{}
	}

	res := make(entity.JoinResponseArray, 0, len(*responses))

	for _, item := range *responses {
		res = append(res, *ToJoinResponse(&item))
	}

	return &res
}

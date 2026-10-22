package view

import (
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

// BuildDepartmentResponse converts a department into its JSON response DTO.
func BuildDepartmentResponse(response *entity.Department) *entity.DepartmentResponse {
	if response == nil {
		return nil
	}

	return &entity.DepartmentResponse{
		ID:          response.ID,
		Name:        response.Name,
		Description: response.Description,
		IsActive:    response.IsActive,
		CreatedAt:   helper.ConvertTimeToString(response.CreatedAt),
		UpdatedAt:   helper.ConvertTimePointerToStringPointer(response.UpdatedAt),
		DeletedAt:   helper.ConvertTimePointerToStringPointer(response.DeletedAt),
	}
}

// BuildDepartmentResponseArray converts departments into response DTOs.
// It returns an empty array when responses is nil or empty.
func BuildDepartmentResponseArray(responses *entity.DepartmentArray) *entity.DepartmentResponseArray {
	if responses == nil || len(*responses) == 0 {
		return &entity.DepartmentResponseArray{}
	}

	res := make(entity.DepartmentResponseArray, 0, len(*responses))

	for _, item := range *responses {
		res = append(res, *BuildDepartmentResponse(&item))
	}

	return &res
}

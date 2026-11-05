package view

import (
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
)

func TestBuildUserResponse(t *testing.T) {
	createdAt := time.Date(2026, time.November, 4, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	deletedAt := createdAt.Add(2 * time.Hour)
	response := BuildUserResponse(&entity.User{
		ID:        "user-1",
		Email:     "user@example.test",
		Password:  "must-not-be-returned",
		IsActive:  true,
		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
		DeletedAt: &deletedAt,
	})

	if response == nil {
		t.Fatal("BuildUserResponse() returned nil")
	}
	if response.ID != "user-1" || response.Email != "user@example.test" || !response.IsActive {
		t.Errorf("response identity = %+v, want the public user fields", response)
	}
	if response.CreatedAt == "" || response.UpdatedAt == nil || response.DeletedAt == nil {
		t.Errorf("response timestamps = %+v, want formatted timestamps", response)
	}
}

func TestBuildUserResponseHandlesNilAndEmptyCollections(t *testing.T) {
	if got := BuildUserResponse(nil); got != nil {
		t.Errorf("BuildUserResponse(nil) = %+v, want nil", got)
	}

	for name, input := range map[string]*[]entity.User{
		"nil slice": nil,
		"empty slice": func() *[]entity.User {
			users := []entity.User{}
			return &users
		}(),
	} {
		t.Run(name, func(t *testing.T) {
			got := BuildUserResponseArray(input)
			if got == nil || len(*got) != 0 {
				t.Fatalf("BuildUserResponseArray() = %v, want a non-nil empty slice", got)
			}
		})
	}
}

func TestBuildUserResponseArrayMapsUsers(t *testing.T) {
	users := []entity.User{
		{ID: "user-1", Email: "one@example.test"},
		{ID: "user-2", Email: "two@example.test"},
	}

	got := BuildUserResponseArray(&users)
	if len(*got) != 2 || (*got)[0].ID != "user-1" || (*got)[1].Email != "two@example.test" {
		t.Fatalf("BuildUserResponseArray() = %+v, want both public user responses", *got)
	}
}

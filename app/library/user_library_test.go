package library

import (
	"errors"
	"net/http"
	"testing"

	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/helper"
)

type userRepositoryStub struct {
	createdUser *entity.User
	user        *entity.User
	users       *[]entity.User
	message     string
	status      int
	err         error
}

func (stub *userRepositoryStub) CreateUser(user *entity.User) (string, int, error) {
	stub.createdUser = user
	return stub.message, stub.status, stub.err
}

func (stub *userRepositoryStub) GetUserByEmail(string) (*entity.User, string, int, error) {
	return stub.user, stub.message, stub.status, stub.err
}

func (stub *userRepositoryStub) GetUsers() (*[]entity.User, string, int, error) {
	return stub.users, stub.message, stub.status, stub.err
}

func TestCreateUserRejectsNilPayload(t *testing.T) {
	lib := NewUserLibrary(nil, &userRepositoryStub{})

	message, status, err := lib.CreateUser(nil)
	if err == nil || status != http.StatusBadRequest || message == "" {
		t.Fatalf("CreateUser(nil) = (%q, %d, %v), want a bad-request error", message, status, err)
	}
}

func TestCreateUserHashesPasswordBeforePersistence(t *testing.T) {
	repository := &userRepositoryStub{message: "created", status: http.StatusOK}
	lib := NewUserLibrary(nil, repository)
	user := &entity.User{ID: "user-1", Email: "user@example.test", Password: "strong-test-password"}

	message, status, err := lib.CreateUser(user)
	if err != nil || status != http.StatusOK || message != "created" {
		t.Fatalf("CreateUser() = (%q, %d, %v), want successful persistence", message, status, err)
	}
	if repository.createdUser != user || user.Password == "strong-test-password" || user.Password == "" {
		t.Fatal("CreateUser() did not pass the user with a hashed password to the repository")
	}
	valid, _, _, err := helper.VerifyArgon("strong-test-password", user.Password)
	if err != nil || !valid {
		t.Fatalf("stored password hash did not verify: valid=%t, err=%v", valid, err)
	}
}

func TestLoginRejectsNilPayload(t *testing.T) {
	lib := NewUserLibrary(nil, &userRepositoryStub{})

	message, status, err := lib.Login(nil, nil)
	if err == nil || status != http.StatusBadRequest || message == "" {
		t.Fatalf("Login(nil) = (%q, %d, %v), want a bad-request error", message, status, err)
	}
}

func TestLoginRejectsIncorrectPassword(t *testing.T) {
	hash, _, _, err := helper.GenerateArgon("correct-password")
	if err != nil {
		t.Fatalf("GenerateArgon() error = %v", err)
	}

	lib := NewUserLibrary(nil, &userRepositoryStub{
		user: &entity.User{Email: "user@example.test", Password: hash, IsActive: true},
	})
	message, status, err := lib.Login(nil, &entity.User{Email: "user@example.test", Password: "incorrect-password"})
	if err == nil || status != http.StatusUnauthorized {
		t.Fatalf("Login(wrong password) = (%q, %d, %v), want unauthorized", message, status, err)
	}
}

func TestGetUsersPropagatesRepositoryResults(t *testing.T) {
	users := []entity.User{{ID: "user-1", Email: "user@example.test"}}
	repository := &userRepositoryStub{users: &users, message: "listed", status: http.StatusOK}
	lib := NewUserLibrary(nil, repository)

	got, message, status, err := lib.GetUsers()
	if err != nil || got != &users || message != "listed" || status != http.StatusOK {
		t.Fatalf("GetUsers() = (%v, %q, %d, %v), want repository result", got, message, status, err)
	}
}

func TestGetUsersWrapsRepositoryErrors(t *testing.T) {
	repositoryErr := errors.New("database unavailable")
	lib := NewUserLibrary(nil, &userRepositoryStub{err: repositoryErr})

	_, _, _, err := lib.GetUsers()
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("GetUsers() error = %v, want wrapped repository error", err)
	}
}

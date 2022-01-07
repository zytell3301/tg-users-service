package core

import (
	"errors"
	"github.com/golang/mock/gomock"
	errors2 "github.com/zytell3301/tg-globals/errors"
	"github.com/zytell3301/tg-users-service/internal/domain"
	"testing"
)

var user = domain.User{
	Name:     "Arshiya",
	Lastname: "Kiani",
	Phone:    "+0000000000",
}

/*
 * Normal case test
 */
func TestService_NewUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	repositoryMock := NewMockUsersRepository(controller)
	repositoryMock.EXPECT().NewUser(user)
	repositoryMock.EXPECT().DoesUserExists(user.Phone)

	core := NewUsersCore(repositoryMock)

	err := core.NewUser(user)

	switch err != nil {
	case true:
		t.Errorf("Expected NewUser to succeed but error returned. Error: %v", err)
	}
}

/*
 * Test case for phone number duplication
 */
func TestService_NewUser2(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	repositoryMock := NewMockUsersRepository(controller)
	repositoryMock.EXPECT().NewUser(user).AnyTimes()
	repositoryMock.EXPECT().DoesUserExists(user.Phone).Return(true, nil)

	core := NewUsersCore(repositoryMock)

	err := core.NewUser(user)
	switch err == nil || !errors.As(err, &UserAlreadyExists{}) {
	case true:
		t.Errorf("Expected NewUser to return error but no error returned")
	}
}

/*
 * Test case for internal failure
 */
func TestService_NewUser3(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	repositoryMock := NewMockUsersRepository(controller)
	repositoryMock.EXPECT().NewUser(user).Return(errors.New(""))
	repositoryMock.EXPECT().DoesUserExists(user.Phone).Return(false, nil)

	core := NewUsersCore(repositoryMock)

	err := core.NewUser(user)
	switch err == nil {
	case true:
		t.Errorf("Expected NewUser to return error but no error returned")
	}

	switch !errors.As(err, &errors2.InternalError{}) {
	case true:
		t.Error("Expected NewUser to return error from type InternalError but error is from different type")
	}
}

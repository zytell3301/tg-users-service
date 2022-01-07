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

var newUsername = "NewUsername"
var dummyError = errors.New("")

func newController(t *testing.T) *gomock.Controller {
	return gomock.NewController(t)
}

/*
 * Normal case test
 */
func TestService_NewUser(t *testing.T) {
	controller := newController(t)
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
	controller := newController(t)
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
	controller := newController(t)
	defer controller.Finish()
	repositoryMock := NewMockUsersRepository(controller)
	repositoryMock.EXPECT().NewUser(user).Return(dummyError)
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

/**
 * Test case for internal failure
 */
func TestService_NewUser4(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := NewMockUsersRepository(controller)
	mock.EXPECT().DoesUserExists(user.Phone).Return(false,dummyError)

	core := NewUsersCore(mock)
	err := core.NewUser(user)
	switch err == nil {
	case true:
		t.Errorf("Expected NewUser to return error but no error returned")
	}
	switch errors.As(err, &errors2.InternalError{}) {
	case false:
		t.Errorf("No proper error returned from NewUser method. Expected NewUser to InternalError error")
	}
}

/**
 * Normal test case
 */
func TestService_UpdateUsername(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := NewMockUsersRepository(controller)
	mock.EXPECT().DoesUsernameExists(newUsername).Return(false, nil)
	mock.EXPECT().UpdateUsername(user.Phone, newUsername)

	core := NewUsersCore(mock)
	err := core.UpdateUsername(user.Phone, newUsername)

	switch err != nil {
	case true:
		t.Errorf("Expected UpdateUsername to succeed but error returned instead. Error message: %v", err)
	}
}

/**
 * Test case for username duplication
 */
func TestService_UpdateUsername2(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := NewMockUsersRepository(controller)
	mock.EXPECT().DoesUsernameExists(newUsername).Return(true, nil)

	core := NewUsersCore(mock)
	err := core.UpdateUsername(user.Phone, newUsername)
	switch err == nil {
	case true:
		t.Errorf("Expected UpdateUsername to return error but no error returned")
	}
	switch errors.As(err, &UsernameAlreadyExists{}) {
	case false:
		t.Errorf("Proper error not returned from UpdateUsername. Expected UpdateUsername to return UsernameAlreadyExists error")
	}
}

/**
 * Test case for internal failure
 */
func TestService_UpdateUsername3(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := NewMockUsersRepository(controller)
	mock.EXPECT().DoesUsernameExists(newUsername).Return(false,dummyError)

	core := NewUsersCore(mock)
	err := core.UpdateUsername(user.Phone,newUsername)
	switch err == nil {
	case true:
		t.Errorf("Expected UpdateUsername to return error but no error returned")
	}
	switch errors.As(err, &errors2.InternalError{}) {
	case false:
		t.Errorf("Proper error not returned from UpdateUsername. Expected UpdateUsername to return UsernameAlreadyExists error")
	}
}

/**
 * Test case for internal failure
 */
func TestService_UpdateUsername4(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := NewMockUsersRepository(controller)
	mock.EXPECT().DoesUsernameExists(newUsername).Return(false,nil)
	mock.EXPECT().UpdateUsername(user.Phone,newUsername).Return(dummyError)

	core := NewUsersCore(mock)
	err := core.UpdateUsername(user.Phone,newUsername)
	switch err == nil {
	case true:
		t.Errorf("Expected UpdateUsername to return error but no error returned")
	}
	switch errors.As(err, &errors2.InternalError{}) {
	case false:
		t.Errorf("Proper error not returned from UpdateUsername. Expected UpdateUsername to return InternalError but no error returned")
	}
}
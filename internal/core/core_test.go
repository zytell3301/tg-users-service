package core

import (
	"bou.ke/monkey"
	"errors"
	"github.com/golang/mock/gomock"
	errors2 "github.com/zytell3301/tg-globals/errors"
	"github.com/zytell3301/tg-users-service/internal/domain"
	"github.com/zytell3301/tg-users-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
)

type qualifyUsername_parameter struct {
	username string
	expected bool
}

var user = domain.User{
	Name:     "Arshiya",
	Lastname: "Kiani",
	Phone:    "+0000000000",
}

var securityCode = domain.SecurityCode{
	Phone: user.Phone,
	Action: security_code_signup_action,
}

var securityCodeRaw = "123456"

var newUsername = "NewUsername"
var dummyError = errors.New("")
var dummyInstanceId = "b8b342e2-3c8a-41f6-8f28-53042ae12519"
var dummyServiceId = "199adc34-f9fd-425e-b721-d5e2b400d289"

func init() {
	hashedSecurityCode, _ := bcrypt.GenerateFromPassword([]byte(securityCodeRaw), 12)
	securityCode.SecurityCode = string(hashedSecurityCode)
}

func newController(t *testing.T) *gomock.Controller {
	return gomock.NewController(t)
}

func patchHasherFunc() {
	monkey.Patch(hashExpression, func(expression string) string {
		return securityCode.SecurityCode
	})
}

/*
 * test case for normal request
 */
func TestService_NewUser(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	repositoryMock := repository.NewMockUsersRepository(controller)
	repositoryMock.EXPECT().GetSecurityCode(user.Phone).Return(securityCode, nil)
	repositoryMock.EXPECT().NewUser(user)
	repositoryMock.EXPECT().DoesUserExists(user.Phone)

	reporterMock := NewMockReporter(controller)

	core := NewUsersCore(repositoryMock, reporterMock, dummyInstanceId, dummyServiceId)

	err := core.NewUser(user, securityCodeRaw)

	switch err != nil {
	case true:
		t.Errorf("Expected NewUser to succeed but error returned. Error: %v Error type: %v", err,reflect.TypeOf(err))
	}
}

/*
 * Test case for phone number duplication
 */
func TestService_NewUser2(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	repositoryMock := repository.NewMockUsersRepository(controller)
	repositoryMock.EXPECT().NewUser(user).AnyTimes()
	repositoryMock.EXPECT().DoesUserExists(user.Phone).Return(true, nil)
	repositoryMock.EXPECT().GetSecurityCode(user.Phone).Return(securityCode, nil)

	reporterMock := NewMockReporter(controller)

	core := NewUsersCore(repositoryMock, reporterMock, dummyInstanceId, dummyServiceId)

	err := core.NewUser(user, securityCodeRaw)
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
	repositoryMock := repository.NewMockUsersRepository(controller)
	repositoryMock.EXPECT().NewUser(user).Return(dummyError)
	repositoryMock.EXPECT().DoesUserExists(user.Phone).Return(false, nil)
	repositoryMock.EXPECT().GetSecurityCode(user.Phone).Return(securityCode, nil)

	reporterMock := NewMockReporter(controller)
	reporterMock.EXPECT().Report(gomock.Any())

	core := NewUsersCore(repositoryMock, reporterMock, dummyInstanceId, dummyServiceId)

	err := core.NewUser(user, securityCodeRaw)
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
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().DoesUserExists(user.Phone).Return(false, dummyError)
	mock.EXPECT().GetSecurityCode(user.Phone).Return(securityCode, nil)

	reporterMock := NewMockReporter(controller)
	reporterMock.EXPECT().Report(gomock.Any())

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
	err := core.NewUser(user, securityCodeRaw)
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
 * test case for normal request
 */
func TestService_UpdateUsername(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().DoesUsernameExists(newUsername).Return(false, nil)
	mock.EXPECT().UpdateUsername(user.Phone, newUsername)

	reporterMock := NewMockReporter(controller)

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
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
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().DoesUsernameExists(newUsername).Return(true, nil)

	reporterMock := NewMockReporter(controller)

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
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
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().DoesUsernameExists(newUsername).Return(false, dummyError)

	reporterMock := NewMockReporter(controller)
	reporterMock.EXPECT().Report(gomock.Any())

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
	err := core.UpdateUsername(user.Phone, newUsername)
	switch err == nil {
	case true:
		t.Errorf("Expected UpdateUsername to return error but no error returned")
	}
	switch errors.As(err, &errors2.InternalError{}) {
	case false:
		t.Errorf("Proper error not returned from UpdateUsername. Expected UpdateUsername to return InternalError error")
	}
}

/**
 * Test case for internal failure
 */
func TestService_UpdateUsername4(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().DoesUsernameExists(newUsername).Return(false, nil)
	mock.EXPECT().UpdateUsername(user.Phone, newUsername).Return(dummyError)

	reporterMock := NewMockReporter(controller)
	reporterMock.EXPECT().Report(gomock.Any())

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
	err := core.UpdateUsername(user.Phone, newUsername)
	switch err == nil {
	case true:
		t.Errorf("Expected UpdateUsername to return error but no error returned")
	}
	switch errors.As(err, &errors2.InternalError{}) {
	case false:
		t.Errorf("Proper error not returned from UpdateUsername. Expected UpdateUsername to return InternalError but no error returned")
	}
}

/**
 * Test case for normal request
 */
func TestService_DeleteUser(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().DeleteUser(user.Phone)

	reporterMock := NewMockReporter(controller)

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
	err := core.DeleteUser(user.Phone)
	switch err != nil {
	case true:
		t.Errorf("Expected DeleteUser to succeed but error returned. Error message: %v", err)
	}
}

/**
 * test case for internal failure
 */
func TestService_DeleteUser2(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().DeleteUser(user.Phone).Return(dummyError)

	reporterMock := NewMockReporter(controller)
	reporterMock.EXPECT().Report(gomock.Any())

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
	err := core.DeleteUser(user.Phone)

	switch err == nil {
	case true:
		t.Errorf("Expected DeleteUser to return error but no error returned")
	}
	switch errors.As(err, &errors2.InternalError{}) {
	case false:
		t.Errorf("Proper error not returned from DeleteUser. Expected DeleteUSer to return InternalError error")
	}
}

/**
 * Normal test case
 */
func TestService_RequestSecurityCode(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().RecordSecurityCode(domain.SecurityCode{
		Phone:        user.Phone,
		SecurityCode: securityCode.SecurityCode,
		Action:       security_code_signup_action,
	}).Return(nil)
	patchHasherFunc()
	defer monkey.UnpatchAll()
	reporterMock := NewMockReporter(controller)

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
	err := core.requestSecurityCode(user.Phone, security_code_signup_action)
	switch err != nil {
	case true:
		t.Errorf("Expected requestSecurityCode method to succeed but an error returned. Error message %v", err)
	}
}

/**
 * Test case for Database failure
 */
func TestService_RequestSecurityCode2(t *testing.T) {
	controller := newController(t)
	defer controller.Finish()
	mock := repository.NewMockUsersRepository(controller)
	mock.EXPECT().RecordSecurityCode(domain.SecurityCode{
		Phone:        user.Phone,
		SecurityCode: securityCode.SecurityCode,
		Action:       security_code_signup_action,
	}).Return(dummyError)

	reporterMock := NewMockReporter(controller)
	reporterMock.EXPECT().Report(gomock.Any())

	core := NewUsersCore(mock, reporterMock, dummyInstanceId, dummyServiceId)
	err := core.requestSecurityCode(user.Phone, security_code_signup_action)
	switch err == nil {
	case true:
		t.Errorf("Expected requestSecurityCode method to return error but no error returned")
	}
}

/**
 * Test cases for invalid parameters
 */
func Test_qualifyUsername(t *testing.T) {
	parameters := []qualifyUsername_parameter{
		{
			// containing space
			username: "aksdjknz adzxc",
			expected: false,
		},
		{
			// starting with number
			username: "1azxsalkc",
			expected: false,
		},
		{
			// too short
			username: "zasd",
			expected: false,
		},
		{
			// starting with number and not containing any character
			username: "124897124",
			expected: false,
		},
		{
			// containing invalid character ( @ )
			username: "@assadklj21",
			expected: false,
		},
		{
			// containing invalid character ( \ )
			username: "\\asdlzkxc",
			expected: false,
		},
		{
			// containing invalid character ( ; )
			username: "xzl;kasdzxasd",
			expected: false,
		},
		{
			// containing invalid character ( ! )
			username: "a!aszcjasd",
			expected: false,
		},
		{
			// containing invalid character ( $ ; )
			username: "$askjzasl;sad",
			expected: false,
		},
		{
			// too long
			username: "ajskzlao1892jkajdkasdhasiodazcasa",
			expected: false,
		},
	}
	for _, parameter := range parameters {
		result := qualifyUsername(parameter.username)
		switch result != parameter.expected {
		case true:
			t.Errorf("Expected qualifyUsername to return false but true returned for invalid username: %v", parameter.username)
		}
	}

}

/**
 * Test case for valid usernames
 */
func Test_qualifyUsername1(t *testing.T) {
	parameters := []qualifyUsername_parameter{
		{
			username: "asklzxcasd",
			expected: true,
		},
		{
			username: "asd1xzcs",
			expected: true,
		},
		{
			username: "a12940841",
			expected: true,
		},
		{
			username: "lkasjdklnzxklcjasdjsklandaskljda",
			expected: true,
		},
	}
	for _, parameter := range parameters {
		result := qualifyUsername(parameter.username)
		switch result != parameter.expected {
		case true:
			t.Errorf("Expected qualifyUsername method to return true but returned false for valid username: %v", parameter.username)
		}
	}
}

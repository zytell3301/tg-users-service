package core

import (
	"github.com/golang/mock/gomock"
	"github.com/zytell3301/tg-users-service/internal/domain"
	"testing"
)

var user = domain.User{
	Name:     "Arshiya",
	Lastname: "Kiani",
	Bio:      "This is an awesome test for bio",
	Username: "aCoolUsername",
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

}

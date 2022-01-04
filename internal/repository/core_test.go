package repository

import (
	"github.com/zytell3301/tg-users-service/internal/domain"
	uuid_generator "github.com/zytell3301/uuid-generator"
	"testing"
	"time"
)

var hosts = []string{"192.168.1.200"}
var idGenerator, _ = uuid_generator.NewGenerator("")
var dummyUser = domain.User{
	Name:          "RK",
	Lastname:      "800",
	Bio:           "This is a test bio",
	Username:      "zimens40K",
	Phone:         "+09999999999",
	Online_status: false,
	Created_at:    time.Now(),
}
var dummyUserId = "5a087beb-4ba5-4583-b2a0-bce500395e1a"
var keyspace = "tg"

func TestNewUsersRepository(t *testing.T) {
	repo, err := NewUsersRepository(hosts, keyspace, idGenerator)
	switch err != nil || repo.connection.Session == nil || repo.connection.Cluster == nil {
	case true:
		t.Errorf("An error encountered while creating a new repo. Error: %v", err)
	}
}

func TestNewUsersRepository2(t *testing.T) {
	_, err := NewUsersRepository(nil, keyspace, idGenerator)
	switch err == nil {
	case true:
		t.Error("Expected to return error but no error returned")
	}
}

// Test fails if the number of current active nodes are less than highest RF (Here it is 3).
// This error is not related to codes
func TestRepository_NewUser(t *testing.T) {
	repo, _ := NewUsersRepository(hosts, keyspace, idGenerator)
	err := repo.NewUser(dummyUser)
	switch err != nil {
	case true:
		t.Errorf("An error encountered while adding a new user. Error: %v", err)
	}
}

func TestRepository_DeleteUser(t *testing.T) {
	repo, _ := NewUsersRepository(hosts, keyspace, idGenerator)
	err := repo.DeleteUser(dummyUser.Phone)
	switch err != nil {
	case true:
		t.Errorf("An error encountered while deleting an existing user. Error: %v", err)
	}
}

func TestRepository_DeleteUser2(t *testing.T) {
	repo, _ := NewUsersRepository(hosts, keyspace, idGenerator)
	err := repo.DeleteUser("")
	switch err == nil {
	case true:
		t.Error("Expected method DeleteUser to return error but no error returned")
	}
}

func TestRepository_UpdateUsername(t *testing.T) {
	repo, _ := NewUsersRepository(hosts, keyspace, idGenerator)
	err := repo.UpdateUsername(dummyUser.Phone, "test_username")
	switch err != nil {
	case true:
		t.Errorf("Expected method UpdateUsername to succeed but error returned. Error message: %v", err)
	}
}

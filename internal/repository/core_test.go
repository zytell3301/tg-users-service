package repository

import (
	"github.com/zytel3301/tg-users-service/internal/domain"
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

func TestNewUsersRepository(t *testing.T) {
	repo, err := NewUsersRepository(hosts, "tg", idGenerator)
	switch err != nil || repo.connection.Session == nil || repo.connection.Cluster == nil {
	case true:
		t.Errorf("An error encountered while creating a new repo. Error: %v", err)
	}
}

func TestNewUsersRepository2(t *testing.T) {
	_, err := NewUsersRepository(nil, "tg", idGenerator)
	switch err == nil {
	case true:
		t.Error("Expected to return error but no error returned")
	}
}

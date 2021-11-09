package repository

import (
	uuid_generator "github.com/zytell3301/uuid-generator"
	"testing"
)

var hosts = []string{"192.168.1.200"}
var idGenerator, _ = uuid_generator.NewGenerator("")

func TestNewUsersRepository(t *testing.T) {
	repo, err := NewUsersRepository(hosts, idGenerator)
	switch err != nil || repo.connection.Session == nil || repo.connection.Cluster == nil {
	case true:
		t.Errorf("An error encountered while creating a new repo. Error: %v", err)
	}
}
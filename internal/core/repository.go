package core

import "github.com/zytell3301/tg-users-service/internal/domain"

type UsersRepository interface {
	NewUser(domain.User) error
	UpdateUsername(string, string) error
	DeleteUser(string) error
	DoesUserExists(string) error
}

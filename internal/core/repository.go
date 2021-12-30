package core

import "github.com/zytell3301/tg-users-service/internal/domain"

type UsersRepository interface {
	NewUser(domain.User) error
	UpdateUser(domain.User) error
	DeleteUser(string) error
}

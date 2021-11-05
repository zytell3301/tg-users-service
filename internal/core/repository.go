package core

import "github.com/zytel3301/tg-users-service/internal/domain"

type UsersRepository interface {
	NewUser(domain.User) error
	UpdateUser(domain.User) error
}

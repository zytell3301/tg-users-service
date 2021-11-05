package core

import "github.com/zytel3301/tg-users-service/internal/domain"

type UsersRepository interface {
	NewUser(user domain.User) error
}

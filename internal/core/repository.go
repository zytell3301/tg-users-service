package core

import "github.com/zytell3301/tg-users-service/internal/domain"

type UsersRepository interface {
	NewUser(user domain.User) error
	UpdateUsername(phone string, username string) error
	DeleteUser(phone string) error
	DoesUserExists(phone string) (bool, error)
	DoesUsernameExists(username string) (bool, error)
	RecordSecurityCode(phone string, code string) error
}

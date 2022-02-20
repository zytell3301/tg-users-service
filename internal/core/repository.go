package core

import "github.com/zytell3301/tg-users-service/internal/domain"

type UsersRepository interface {
	NewUser(user domain.User) error
	UpdateUsername(phone string, username string) error
	DeleteUser(phone string) error
	DoesUserExists(phone string) (bool, error)
	DoesUsernameExists(username string) (bool, error)
	RecordSecurityCode(securityCode domain.SecurityCode) error
	GetSecurityCode(phone string) (domain.SecurityCode, error)
	GetUserByPhone(phone string) (domain.User, error)
	GetUserByUsername(username string) (domain.User, error)
}

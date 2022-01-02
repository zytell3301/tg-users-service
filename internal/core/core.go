package core

import "github.com/zytell3301/tg-users-service/internal/domain"

type Service struct {
	repository UsersRepository
}

func NewUsersCore(repository UsersRepository) Service {
	return Service{
		repository: repository,
	}
}

// @TODO an error must be returned if an error encountered while checking for existing user. A global package must be created for this purpose
func (s Service) NewUser(user domain.User) (err error) {
	doesExists, err := s.repository.DoesUserExists(user.Phone)
	switch doesExists {
	case true:
		return UserAlreadyExists{}
	}
	err = s.repository.NewUser(user)
	switch err != nil {
	case true:
		// @TODO once the logger service implemented, this part must report the error to logger service
	}

	return
}

// @TODO check if the username already exists or not
func (s Service) UpdateUsername(phone string, username string) (err error) {
	err = s.repository.UpdateUsername(phone, username)
	switch err != nil {
	case true:
		// @TODO once the logger service implemented, this part must report the error to logger service
	}

	return
}

func (s Service) DeleteUser(phone string) (err error) {
	err = s.repository.DeleteUser(phone)
	switch err != nil {
	case true:
		// @TODO once the logger service implemented, this part must report the error to logger service
	}

	return
}

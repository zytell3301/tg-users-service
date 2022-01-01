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

// @TODO it must be checked whether the user already exists or not
func (s Service) NewUser(user domain.User) (err error) {
	err = s.repository.NewUser(user)
	switch err != nil {
	case true:
		// @TODO once the logger service implemented, this part must report the error to logger service
	}

	return
}

// @TODO check if the username already exists or not
func (s Service) UpdateUsername(user domain.User, username string) (err error) {
	err = s.repository.UpdateUsername(user.Phone, username)
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

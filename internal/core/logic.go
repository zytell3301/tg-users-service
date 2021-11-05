package core

import "github.com/zytel3301/tg-users-service/internal/domain"

type Service struct {
	repository UsersRepository
}

func NewUsersCore(repository UsersRepository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) NewUser(user domain.User) (err error) {
	err = s.repository.NewUser(user)
	switch err != nil {
	case true:
		// @TODO once the logger service implemented, this part must report the error to logger service
	}

	return
}

func (s Service) UpdateUser(user domain.User) (err error) {
	err = s.repository.UpdateUser(user)
	switch err != nil {
	case true:
		// @TODO once the logger service implemented, this part must report the error to logger service
	}

	return
}

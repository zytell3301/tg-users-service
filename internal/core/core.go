package core

import (
	"github.com/zytell3301/tg-globals/errors"
	"github.com/zytell3301/tg-users-service/internal/domain"
)

type Service struct {
	repository UsersRepository
}

func NewUsersCore(repository UsersRepository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) NewUser(user domain.User) (err error) {
	doesExists, err := s.repository.DoesUserExists(user.Phone)
	switch err != nil {
	case true:
		// @TODO report error to central error recorder
		return errors.InternalError{}
	}
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

func (s Service) UpdateUsername(phone string, username string) (err error) {
	doesExists, err := s.repository.DoesUsernameExists(username)
	switch err != nil {
	case true:
		// @TODO error must be reported to central error recorder
		return errors.InternalError{}
	}
	switch doesExists {
	case true:
		return UsernameAlreadyExists{}
	}
	err = s.repository.UpdateUsername(phone, username)
	switch err != nil {
	case true:
		// @TODO once the logger service implemented, this part must report the error to logger service
		return errors.InternalError{}
	}

	return
}

func (s Service) DeleteUser(phone string) (err error) {
	err = s.repository.DeleteUser(phone)
	switch err != nil {
	case true:
		// @TODO once the logger service implemented, this part must report the error to logger service
		return errors.InternalError{}
	}

	return
}

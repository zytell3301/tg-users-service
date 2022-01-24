package core

import (
	"crypto/rand"
	"github.com/zytell3301/tg-error-reporter"
	"github.com/zytell3301/tg-globals/errors"
	"github.com/zytell3301/tg-users-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"regexp"
	"strconv"
)

type Service struct {
	repository    UsersRepository
	ErrorReporter ErrorReporter.Reporter
	instanceId    string
	serviceId     string
}

func NewUsersCore(repository UsersRepository, errorReporter ErrorReporter.Reporter, instanceId string, serviceId string) Service {
	return Service{
		repository:    repository,
		ErrorReporter: errorReporter,
		instanceId:    instanceId,
		serviceId:     serviceId,
	}
}

func (s Service) NewUser(user domain.User) (err error) {
	doesExists, err := s.repository.DoesUserExists(user.Phone)
	switch err != nil {
	case true:
		s.ErrorReporter.Report(ErrorReporter.Error{
			ServiceId:  s.serviceId,
			InstanceId: s.instanceId,
			Message:    err.Error(),
		})
		return errors.InternalError{}
	}
	switch doesExists {
	case true:
		return UserAlreadyExists{}
	}
	err = s.repository.NewUser(domain.User{
		Name:     user.Name,
		Lastname: user.Lastname,
		Phone:    user.Phone,
	})
	switch err != nil {
	case true:
		s.ErrorReporter.Report(ErrorReporter.Error{
			ServiceId:  s.serviceId,
			InstanceId: s.instanceId,
			Message:    err.Error(),
		})
		return errors.InternalError{}
	}

	return
}

// @TODO qualify username before processing request
func (s Service) UpdateUsername(phone string, username string) (err error) {
	doesExists, err := s.repository.DoesUsernameExists(username)
	switch err != nil {
	case true:
		s.ErrorReporter.Report(ErrorReporter.Error{
			ServiceId:  s.serviceId,
			InstanceId: s.instanceId,
			Message:    err.Error(),
		})
		return errors.InternalError{}
	}
	switch doesExists {
	case true:
		return UsernameAlreadyExists{}
	}
	err = s.repository.UpdateUsername(phone, username)
	switch err != nil {
	case true:
		s.ErrorReporter.Report(ErrorReporter.Error{
			ServiceId:  s.serviceId,
			InstanceId: s.instanceId,
			Message:    err.Error(),
		})
		return errors.InternalError{}
	}

	return
}

func qualifyUsername(username string) bool {
	isValid, _ := regexp.MatchString("^\\D[\\w,\\d,_]{7,31}$", username)
	return isValid
}

func (s Service) DeleteUser(phone string) (err error) {
	err = s.repository.DeleteUser(phone)
	switch err != nil {
	case true:
		s.ErrorReporter.Report(ErrorReporter.Error{
			ServiceId:  s.serviceId,
			InstanceId: s.instanceId,
			Message:    err.Error(),
		})
		return errors.InternalError{}
	}

	return
}

func (s Service) RequestSecurityCode(phone string) (err error) {
	securityCode := hashExpression(generateSecurityCode())
	err = s.repository.RecordSecurityCode(phone, securityCode)
	switch err != nil {
	case true:
		s.ErrorReporter.Report(ErrorReporter.Error{
			ServiceId:  s.serviceId,
			InstanceId: s.instanceId,
			Message:    err.Error(),
		})
	}
	return
}

func hashExpression(expression string) string {
	hashedExpression, _ := bcrypt.GenerateFromPassword([]byte(expression), 12)
	return string(hashedExpression)
}

func generateSecurityCode() string {
	securityCode, _ := rand.Int(rand.Reader, big.NewInt(999999-100000))
	return strconv.Itoa(int(securityCode.Int64()) + 100000)
}

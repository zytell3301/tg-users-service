package core

import (
	"crypto/rand"
	"fmt"
	"github.com/zytell3301/tg-error-reporter"
	"github.com/zytell3301/tg-globals/errors"
	"github.com/zytell3301/tg-users-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
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

func (s Service) UpdateUsername(phone string, username string) (err error) {
	switch qualifyUsername(username) {
	case false:
		return UsernameNotQualifiedError
	}
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
	invalidCharacters := []string{
		"@",
		"\\",
	}
	for _, invalidCharacter := range invalidCharacters {
		switch strings.Contains(username, invalidCharacter) {
		case true:
			return false
		}
	}
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

func (s Service) VerifySecurityCode(phone string, code string) error {
	securityCode, err := s.repository.GetSecurityCode(phone)
	switch err != nil {
	case true:
		s.ErrorReporter.Report(ErrorReporter.Error{
			ServiceId:  s.serviceId,
			InstanceId: s.instanceId,
			Message:    fmt.Sprintf("An error occurred while fetching security code from repository. Error message: %s", err.Error()),
		})
	}
	switch checkHashMatch(code, securityCode.SecurityCode) {
	case false:
		return SecurityCodeNotValidError
	}
	return nil
}

func hashExpression(expression string) string {
	hashedExpression, _ := bcrypt.GenerateFromPassword([]byte(expression), 12)
	return string(hashedExpression)
}

func generateSecurityCode() string {
	securityCode, _ := rand.Int(rand.Reader, big.NewInt(999999-100000))
	return strconv.Itoa(int(securityCode.Int64()) + 100000)
}

func checkHashMatch(expression string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(expression)) == nil
}

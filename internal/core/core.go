package core

import (
	"crypto/rand"
	errors2 "errors"
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

const (
	security_code_signup_action = "SIGN UP"
)

func NewUsersCore(repository UsersRepository, errorReporter ErrorReporter.Reporter, instanceId string, serviceId string) Service {
	return Service{
		repository:    repository,
		ErrorReporter: errorReporter,
		instanceId:    instanceId,
		serviceId:     serviceId,
	}
}

func (s Service) NewUser(user domain.User, securityCode string) (err error) {
	err = s.VerifySecurityCode(user.Phone, securityCode, security_code_signup_action)
	switch err != nil {
	case true:
		switch errors2.As(err, &SecurityCodeNotValid{}) {
		case true:
			return err
		default:
			s.ErrorReporter.Report(ErrorReporter.Error{
				ServiceId:  s.serviceId,
				InstanceId: s.instanceId,
				Message:    fmt.Sprintf("An error occurred while checking for security code. Error message: %s", err.Error()),
			})
			return errors.InternalErrorOccurred
		}
	}
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

func (s Service) RequestSignupSecurityCode(phone string) error {
	doesExists, err := s.repository.DoesUserExists(phone)
	switch err != nil {
	case true:
		s.ErrorReporter.Report(ErrorReporter.Error{
			ServiceId:  s.serviceId,
			InstanceId: s.instanceId,
			Message:    fmt.Sprintf("An error occurred while checking for user existence. Error message: %v", err.Error()),
		})
		return errors.InternalErrorOccurred
	}
	switch doesExists {
	case true:
		return UserAlreadyExistsError
	default:
		return s.RequestSecurityCode(phone, security_code_signup_action)
	}
}

func (s Service) RequestSecurityCode(phone string, action string) (err error) {
	securityCode := hashExpression(generateSecurityCode())
	err = s.repository.RecordSecurityCode(domain.SecurityCode{
		Phone:        phone,
		Action:       action,
		SecurityCode: securityCode,
	})
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

func (s Service) VerifySecurityCode(phone string, code string, action string) error {
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
	switch securityCode.Action != action {
	case true:
		return SecurityCodeActionDoesNotMatchError
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

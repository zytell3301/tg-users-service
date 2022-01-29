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
	security_code_login_action  = "LOGIN"
)

func NewUsersCore(repository UsersRepository, errorReporter ErrorReporter.Reporter, instanceId string, serviceId string) Service {
	return Service{
		repository:    repository,
		ErrorReporter: errorReporter,
		instanceId:    instanceId,
		serviceId:     serviceId,
	}
}

/**
 * Creates a new user if the phone number already exists. Otherwise it returns UserAlreadyExists error
 */
func (s Service) NewUser(user domain.User, securityCode string) (err error) {
	err = s.VerifySecurityCode(user.Phone, securityCode, security_code_signup_action)
	switch err != nil {
	case true:
		switch errors2.As(err, &SecurityCodeNotValid{}) {
		case true:
			return err
		default:
			return errors.InternalError{}
		}
	}
	doesExists, err := s.repository.DoesUserExists(user.Phone)
	switch err != nil {
	case true:
		s.reportDoesUserExistsError(err)
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
		s.reportNewUserError(err)
		return errors.InternalError{}
	}

	return
}

/**
 * Updates current user's username or sets a new one if the user currently don't have username.
 * First username is qualified under username policies and then the username existence is checked before update.
 * If the username qualification failed UsernameNotQualified error is returned.
 * If the username exists UsernameAlreadyExists error will be returned.
 */
func (s Service) UpdateUsername(phone string, username string) (err error) {
	switch qualifyUsername(username) {
	case false:
		return UsernameNotQualified{}
	}
	doesExists, err := s.repository.DoesUsernameExists(username)
	switch err != nil {
	case true:
		s.reportDoesUsernameExistsError(err)
		return errors.InternalError{}
	}
	switch doesExists {
	case true:
		return UsernameAlreadyExists{}
	}
	err = s.repository.UpdateUsername(phone, username)
	switch err != nil {
	case true:
		s.reportUpdateUsernameError(err)
		return errors.InternalError{}
	}

	return
}

/**
 * Username qualification rules:
 * least username length is 8
 * max username length is 32
 * username must only contain english characters, digits and underscore( _ )
 */
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

/**
 * Deletes user account.
 * @TODO other user data must be deleted like messages
 */
func (s Service) DeleteUser(phone string) (err error) {
	err = s.repository.DeleteUser(phone)
	switch err != nil {
	case true:
		s.reportDeleteUserError(err)
		return errors.InternalError{}
	}

	return
}

/**
 * Creates a new security code for only signing up.
 * If the user already exists UserAlreadyExists error will be returned.
 */
func (s Service) RequestSignupSecurityCode(phone string) error {
	doesExists, err := s.repository.DoesUserExists(phone)
	switch err != nil {
	case true:
		s.reportDoesUserExistsError(err)
		return errors.InternalError{}
	}
	switch doesExists {
	case true:
		return UserAlreadyExists{}
	default:
		return s.requestSecurityCode(phone, security_code_signup_action)
	}
}

/**
 * Creates a new security code only for login.
 * If the user does not exists UserNotFound error will be returned
 */
func (s Service) RequestLoginSecurityCode(phone string) error {
	doesExists, err := s.repository.DoesUserExists(phone)
	switch err != nil {
	case true:
		s.reportDoesUserExistsError(err)
	}
	switch doesExists {
	case false:
		return UserNotFound{}
	default:
		return s.requestSecurityCode(phone, security_code_login_action)
	}
}

/**
 * Creates a new security code but this method is not directly accessible from outside of package.
 * It is only available from RequestLoginSecurityCode or RequestSignupSecurityCode methods
 */
func (s Service) requestSecurityCode(phone string, action string) (err error) {
	securityCode := hashExpression(generateSecurityCode())
	err = s.repository.RecordSecurityCode(domain.SecurityCode{
		Phone:        phone,
		Action:       action,
		SecurityCode: securityCode,
	})
	switch err != nil {
	case true:
		s.reportRequestSecurityCodeError(err)
	}
	return
}

/**
 * Verifies given security code and action.
 * If the security code is incorrect SecurityCodeNotValid error will be returned.
 * If the security code is correct but the action is incorrect, SecurityCodeActionDoesNotMatch will be returned
 */
func (s Service) VerifySecurityCode(phone string, code string, action string) error {
	securityCode, err := s.repository.GetSecurityCode(phone)
	switch err != nil {
	case true:
		s.reportGetSecurityCodeError(err)
		return errors.InternalError{}
	}
	switch checkHashMatch(code, securityCode.SecurityCode) {
	case false:
		return SecurityCodeNotValid{}
	}
	switch securityCode.Action != action {
	case true:
		return SecurityCodeActionDoesNotMatch{}
	}
	return nil
}

/**
 * Reports service errors to central error recorder.
 * Message layout can be in format of fmt.Sprintf with its parameters
 */
func (s Service) reportError(message string, parameters ...string) {
	go s.ErrorReporter.Report(ErrorReporter.Error{
		ServiceId:  s.serviceId,
		InstanceId: s.instanceId,
		Message:    fmt.Sprintf(message, parameters),
	})
}

/**
 * Reports service errors to central error recorder with a pre-defined message
 */
func (s Service) reportErr(subject string, err error) {
	s.reportError("An error occurred while %s. Error message: %s", subject, err.Error())
}

func (s Service) reportGetSecurityCodeError(err error) {
	s.reportErr("fetching security code from repository", err)
}

func (s Service) reportDoesUserExistsError(err error) {
	s.reportErr("checking for user existence", err)
}

func (s Service) reportDeleteUserError(err error) {
	s.reportErr("deleting user from database", err)
}

func (s Service) reportRequestSecurityCodeError(err error) {
	s.reportErr("recording security code on database", err)
}

func (s Service) reportDoesUsernameExistsError(err error) {
	s.reportErr("checking for username existence", err)
}

func (s Service) reportUpdateUsernameError(err error) {
	s.reportErr("updating username in database", err)
}

func (s Service) reportNewUserError(err error) {
	s.reportErr("inserting user into database", err)
}

func hashExpression(expression string) string {
	hashedExpression, _ := bcrypt.GenerateFromPassword([]byte(expression), 12)
	return string(hashedExpression)
}

func generateSecurityCode() string {
	securityCode, _ := rand.Int(rand.Reader, big.NewInt(999999-100000))
	return strconv.Itoa(int(securityCode.Int64()) + 100000)
}

// @TODO check for hash cost and swap with correct one if the cost does not match
func checkHashMatch(expression string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(expression)) == nil
}

package core

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	errors2 "errors"
	"github.com/zytell3301/tg-globals/errors"
	"github.com/zytell3301/tg-users-service/internal/domain"
	"github.com/zytell3301/tg-users-service/internal/errorReporter"
	"github.com/zytell3301/tg-users-service/pkg/CertGen"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

type Service struct {
	repository    UsersRepository
	certGen       CertGen.Gen
}

const (
	security_code_signup_action = "SIGN UP"
	security_code_login_action  = "LOGIN"
)

func NewUsersCore(repository UsersRepository, certGen CertGen.Gen) Service {
	return Service{
		repository: repository,
		certGen:    certGen,
	}
}

/**
 * Creates a new user if the phone number already exists. Otherwise it returns UserAlreadyExists error
 */
func (s Service) NewUser(user domain.User, securityCode string) (err error) {
	err = s.VerifySecurityCode(user.Phone, securityCode, security_code_signup_action)
	switch err != nil {
	case true:
		switch errors2.As(err, &errors.InternalError{}) {
		case true:
			return errors.InternalError{}
		default:
			return err
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
 * Generate a certificate for corresponding user if provided security code is correct.
 */
func (s Service) Login(phone string, securityCode string) ([]byte, error) {
	err := s.VerifySecurityCode(phone, securityCode, security_code_login_action)
	switch err != nil {
	case true:
		switch errors2.As(err, &SecurityCodeNotValid{}) {
		case true:
			return nil, err
		default:
			return nil, errors.InternalError{}
		}
	}
	user, err := s.repository.GetUserByPhone(phone)
	switch err != nil {
	case true:
		switch errors2.As(err, errors.EntityNotFound{}) {
		case true:
			return nil, UserNotFound{}
		default:
			s.reportGetUserByPhoneError(err)
			return nil, errors.InternalError{}
		}
	}
	cert, err := s.generateUserCert(user)
	switch err != nil {
	case true:
		return nil, errors.InternalError{}
	}
	return cert, nil
}

/**
 * Generates a certificate based on user credentials
 */
func (s Service) generateUserCert(user domain.User) ([]byte, error) {
	cert, err := s.certGen.NewCertificate(&x509.Certificate{
		Subject: pkix.Name{
			SerialNumber: user.Id,
		},
	})
	switch err != nil {
	case true:
		s.reportError("generating certificate", err)
	}
	return cert, err
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

func (s Service) GetUserByUsername(username string) (domain.User, error) {
	user, err := s.repository.GetUserByUsername(username)
	switch err != nil {
	case true:
		switch errors2.As(err, errors.EntityNotFound{}) {
		case true:
			return domain.User{}, UserNotFound{}
		}
		return domain.User{}, errors.InternalError{}
	}
	return user, nil
}

/**
 * Reports service errors to central error recorder with a pre-defined message
 */
func (s Service) reportError(subject string, err error) {
	errorReporter.ReportError("An error occurred while %s. Error message: %s", subject, err.Error())
}

func (s Service) reportGetSecurityCodeError(err error) {
	s.reportError("fetching security code from repository", err)
}

func (s Service) reportDoesUserExistsError(err error) {
	s.reportError("checking for user existence", err)
}

func (s Service) reportDeleteUserError(err error) {
	s.reportError("deleting user from database", err)
}

func (s Service) reportRequestSecurityCodeError(err error) {
	s.reportError("recording security code on database", err)
}

func (s Service) reportDoesUsernameExistsError(err error) {
	s.reportError("checking for username existence", err)
}

func (s Service) reportUpdateUsernameError(err error) {
	s.reportError("updating username in database", err)
}

func (s Service) reportNewUserError(err error) {
	s.reportError("inserting user into database", err)
}

func (s Service) reportGetUserByPhoneError(err error) {
	s.reportError("fetching user from database", err)
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

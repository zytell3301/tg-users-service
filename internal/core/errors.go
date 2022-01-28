package core

import (
	"github.com/zytell3301/tg-globals/errors"
)

type UserAlreadyExists struct {
	errors.Derror
}

type UserNotFound struct {
	errors.Derror
}

type UsernameAlreadyExists struct {
	errors.Derror
}

type UsernameNotQualified struct {
	errors.Derror
}

type SecurityCodeNotValid struct {
	errors.Derror
}

type SecurityCodeActionDoesNotMatch struct {
	errors.Derror
}

var (
	UserAlreadyExistsError = UserAlreadyExists{
		errors.Derror{
			Message: "user already exists",
			Code:    2,
		},
	}
	UserNotFoundError = UserNotFound{
		errors.Derror{
			Message: "user not found",
			Code:    3,
		},
	}
	UsernameAlreadyExistsError = UsernameAlreadyExists{
		errors.Derror{
			Message: "username already exists",
			Code:    4,
		},
	}
	UsernameNotQualifiedError = UsernameNotQualified{
		errors.Derror{
			Message: "username not qualified",
			Code:    6,
		},
	}
	SecurityCodeNotValidError = SecurityCodeNotValid{
		errors.Derror{
			Message: "security code is incorrect or expired ",
			Code:    7,
		},
	}
	SecurityCodeActionDoesNotMatchError = SecurityCodeActionDoesNotMatch{
		errors.Derror{
			Message: "security code action does not match with current action",
			Code:    8,
		},
	}
)

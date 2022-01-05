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

type UsernameTooShort struct {
	errors.Derror
}

var (
	UserAlreadyExistsError = UserAlreadyExists{
		errors.Derror{
			Message: "user already exists",
			Code:    1,
		},
	}
	UserNotFoundError = UserNotFound{
		errors.Derror{
			Message: "user not found",
			Code:    1,
		},
	}
	UsernameAlreadyExistsError = UsernameAlreadyExists{
		errors.Derror{
			Message: "username already exists",
			Code:    1,
		},
	}
	UsernameTooShortError = UsernameTooShort{
		errors.Derror{
			Message: "username too short",
			Code:    2,
		},
	}
)

package core

type derror struct {
	Message string
	Code    uint32
}

type UserAlreadyExists struct {
	derror
}

type UserNotFound struct {
	derror
}

type UsernameAlreadyExists struct {
	derror
}

type UsernameTooShort struct {
	derror
}

var (
	UserAlreadyExistsError = UserAlreadyExists{
		derror{
			Message: "user already exists",
			Code:    1,
		},
	}
	UserNotFoundError = UserNotFound{
		derror{
			Message: "user not found",
			Code:    1,
		},
	}
	UsernameAlreadyExistsError = UsernameAlreadyExists{
		derror{
			Message: "username already exists",
			Code:    1,
		},
	}
	UsernameTooShortError = UsernameTooShort{
		derror{
			Message: "username too short",
			Code:    2,
		},
	}
)

func (e derror) Error() string {
	return e.Message
}

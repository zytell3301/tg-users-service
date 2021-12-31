package core

type derror struct {
	Message string
	Code    uint32
}

type UserAlreadyExists struct {
	derror
}

var (
	UserAlreadyExistsError = UserAlreadyExists{
		derror{
			Message: "user already exists",
			Code:    1,
		},
	}
)

func (e derror) Error() string {
	return e.Message
}

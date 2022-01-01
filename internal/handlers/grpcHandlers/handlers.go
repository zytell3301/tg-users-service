package grpcHandlers

import (
	"context"
	"errors"
	"github.com/zytell3301/tg-users-service/internal/core"
	"github.com/zytell3301/tg-users-service/internal/domain"
	"github.com/zytell3301/tg-users-service/pkg/UsersService"
	error1 "github.com/zytell3301/tg-users-service/pkg/error"
	"time"
)

type Handler struct {
	UsersService.UnimplementedUsersServiceServer
	core core.Service
}

/**
 * error core 0 indicates that the operation completed successfully with no error in all handlers
 */

func (h Handler) NewUser(ctx context.Context, user *UsersService.User) (*error1.Error, error) {
	err := h.core.NewUser(domain.User{
		Name:       user.Name,
		Lastname:   user.Lastname,
		Phone:      user.Phone,
		Created_at: time.Now(),
	})

	switch {
	case errors.As(err, &core.UserAlreadyExists{}):
		return &error1.Error{
			Message: core.UserAlreadyExistsError.Message,
			Code:    core.UserAlreadyExistsError.Code,
		}, nil
	}

	return &error1.Error{
		Code: 0,
	}, nil
}

func (h Handler) DeleteUser(ctx context.Context, phone *UsersService.Phone) (*error1.Error, error) {
	err := h.core.DeleteUser(phone.Phone)

	switch {
	case errors.As(err, core.UserNotFound{}):
		return &error1.Error{
			Message: core.UserNotFoundError.Message,
			Code:    core.UserNotFoundError.Code,
		}, nil
	}

	return &error1.Error{
		Code: 0,
	}, nil
}

func (h Handler) UpdateUsername(ctx context.Context, message *UsersService.UpdateUsernameMessage) (*error1.Error, error) {
	err := h.core.UpdateUsername(message.Phone, message.Username)

	switch {
	case errors.As(err, core.UsernameTooShort{}):
		return &error1.Error{
			Message: core.UsernameTooShortError.Message,
			Code:    core.UsernameTooShortError.Code,
		}, nil

	case errors.As(err, core.UsernameAlreadyExists{}):
		return &error1.Error{
			Message: core.UsernameAlreadyExistsError.Message,
			Code:    core.UsernameAlreadyExistsError.Code,
		}, nil
	}

	return &error1.Error{
		Code: 0,
	}, nil
}

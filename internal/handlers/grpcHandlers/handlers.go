package grpcHandlers

import (
	"context"
	"errors"
	errors2 "github.com/zytell3301/tg-globals/errors"
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

func NewHandler(service core.Service) Handler {
	return Handler{
		core: service,
	}
}

func (h Handler) NewUser(ctx context.Context, message *UsersService.NewUserMessage) (*error1.Error, error) {
	err := h.core.NewUser(domain.User{
		Name:       message.User.Name,
		Lastname:   message.User.Lastname,
		Phone:      message.User.Phone,
		Created_at: time.Now(),
	}, message.SecurityCode.Code)

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
	case errors.As(err, &core.UserNotFound{}):
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
	case errors.As(err, &core.UsernameNotQualified{}):
		return &error1.Error{
			Message: core.UsernameNotQualifiedError.Message,
			Code:    core.UsernameNotQualifiedError.Code,
		}, nil

	case errors.As(err, &core.UsernameAlreadyExists{}):
		return &error1.Error{
			Message: core.UsernameAlreadyExistsError.Message,
			Code:    core.UsernameAlreadyExistsError.Code,
		}, nil
	}

	return &error1.Error{
		Code: 0,
	}, nil
}

func (h Handler) Login(_ context.Context, request *UsersService.LoginRequest) (*UsersService.LoginResponse, error) {
	cert, err := h.core.Login(request.Phone, request.Phone)
	switch {
	case errors.As(err, &core.SecurityCodeNotValid{}):
		return &UsersService.LoginResponse{
			Error: &error1.Error{
				Message: core.SecurityCodeNotValidError.Message,
				Code:    core.SecurityCodeNotValidError.Code,
			},
		}, nil
	case errors.As(err, &errors2.InternalError{}):
		return &UsersService.LoginResponse{
			Error: &error1.Error{
				Message: errors2.InternalErrorOccurred.Message,
				Code:    errors2.InternalErrorOccurred.Code,
			},
		}, nil
	case errors.As(err, &core.UserNotFound{}):
		return &UsersService.LoginResponse{
			Error: &error1.Error{
				Message: core.UserNotFoundError.Message,
				Code:    core.UserNotFoundError.Code,
			},
		}, nil
	}
	return &UsersService.LoginResponse{
		Certificate: cert,
	}, nil
}

func (h Handler) RequestSignupSecurityCode(_ context.Context, request *UsersService.Phone) (*error1.Error, error) {
	err := h.core.RequestSignupSecurityCode(request.Phone)
	switch {
	case errors.As(err, &errors2.InternalError{}):
		return &error1.Error{
			Message: errors2.InternalErrorOccurred.Message,
			Code:    errors2.InternalErrorOccurred.Code,
		}, nil
	case errors.As(err, &core.UserAlreadyExists{}):
		return &error1.Error{
			Message: core.UserAlreadyExistsError.Message,
			Code:    core.UserAlreadyExistsError.Code,
		}, nil
	}
	return nil, nil
}

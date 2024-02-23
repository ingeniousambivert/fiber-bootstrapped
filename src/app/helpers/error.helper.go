package helpers

import (
	"errors"

	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/utils"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func NotFound(m string) *core.ServerError {
	return &core.ServerError{Status: utils.HttpStatusNotFound, Title: "not-found", Message: m}
}

func BadRequest(m string) *core.ServerError {
	return &core.ServerError{Status: utils.HttpStatusBadRequest, Title: "bad-request", Message: m}
}

func Unauthorized(m string) *core.ServerError {
	return &core.ServerError{Status: utils.HttpStatusUnauthorized, Title: "unauthorized", Message: m}
}
func Forbidden(m string) *core.ServerError {
	return &core.ServerError{Status: utils.HttpStatusForbidden, Title: "forbidden", Message: m}
}
func Conflict(m string) *core.ServerError {
	return &core.ServerError{Status: utils.HttpStatusConflict, Title: "conflict", Message: m}
}

func Unexpected(m string) *core.ServerError {
	return &core.ServerError{Status: utils.HttpStatusInternalServerError, Title: "internal-server", Message: m}
}

func CreateError(message string) error {
	return errors.New(message)
}

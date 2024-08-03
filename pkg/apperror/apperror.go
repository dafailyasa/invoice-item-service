package apperror

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrCustomerExistInvoice = errors.New("can't delete customer. because have invoices")

	ErrCustomerEmailExist = errors.New("customer email already used")
	ErrItemNameExist      = errors.New("item name already used")

	ErrInvNotFound      = errors.New("invoice was not found")
	ErrCustomerNotFound = errors.New("customer was not found")
	ErrItemNotFound     = errors.New("item was not found")
)

type AppError struct {
	Code    int
	Err     error
	Message string
}

func (h AppError) Error() string {
	return h.Err.Error()
}

func BadRequest(err error) error {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "bad_request",
		Err:     err,
	}
}

func UnprocessableEntity(err error) error {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "unprocessable_entity",
		Err:     err,
	}
}

func InternalServerError(err error) error {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal_server_error",
		Err:     err,
	}
}

func Unauthorized(err error) error {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
		Err:     err,
	}
}

func Forbidden(err error) error {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: "forbidden",
		Err:     err,
	}
}

func NotFound(err error) error {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: "not_found",
		Err:     err,
	}
}

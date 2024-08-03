package response

import (
	"errors"
	"net/http"

	"github.com/dafailyasa/invoice-item-service/pkg/apperror"
	"github.com/labstack/echo/v4"
)

const (
	InternalServerError = "internal_server_error"
	SUCCESS             = "success"
)

// FailedResponse represents a failed response structure for API responses.
type FailedResponse struct {
	Code    int    `json:"code" `   // HTTP status code.
	Message string `json:"message"` // Message corresponding to the status code.
	Error   string `json:"error"`   // error message.
}

// ErrorBuilder constructs a FailedResponse based on the provided error.
func ErrorBuilder(err error) FailedResponse {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		ae := err.(*apperror.AppError)

		return FailedResponse{
			Code:    ae.Code,
			Message: ae.Message,
			Error:   ae.Error(),
		}
	}

	var errString = InternalServerError
	if err != nil {
		errString = err.Error()
	}

	return FailedResponse{
		Code:    http.StatusInternalServerError,
		Message: InternalServerError,
		Error:   errString,
	}
}

func (x FailedResponse) Send(c echo.Context) error {
	return c.JSON(x.Code, x)
}

// SuccessResponse represents a success response structure for API responses.
type SuccessResponse struct {
	Success
	Pagination
}

type ResponseFormat struct {
	Code    int    `json:"code"` // HTTP status code.
	Message string `json:"message"`
}

type Success struct {
	ResponseFormat
	Data interface{} `json:"data,omitempty"` // data payload.
}

type Pagination struct {
	Pagination interface{} `json:"pagination,omitempty"` //pagination payload.
	Success
}

func SuccessBuilder(response interface{}, pagination ...interface{}) SuccessResponse {
	result := SuccessResponse{
		Success: Success{
			ResponseFormat: ResponseFormat{
				Code:    http.StatusOK,
				Message: SUCCESS,
			},
			Data: response,
		},
	}

	if len(pagination) > 0 {
		result.Pagination.Pagination = pagination[0]
	}

	return result
}

func (c SuccessResponse) Send(ctx echo.Context) error {
	return ctx.JSON(c.Code, c)
}

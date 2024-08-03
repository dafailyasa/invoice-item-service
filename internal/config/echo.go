package config

import (
	"errors"
	"net/http"

	"github.com/dafailyasa/invoice-item-service/pkg/apperror"
	"github.com/dafailyasa/invoice-item-service/pkg/response"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// configCORS contains the CORS (Cross-Origin Resource Sharing) configuration for the server.
var configCORS = echoMiddleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
}

// NewEchoServer creates and configures a new Echo server instance.
// Parameters:
//   - cfg: The application configuration.
//
// Returns:
//   - *echo.Echo: A configured Echo server instance.
func NewEchoServer(cfg Config) *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{DisableStackAll: true}))
	e.Use(echoMiddleware.CORSWithConfig(configCORS))

	e.HTTPErrorHandler = errorHandler
	e.Debug = cfg.Server.Debug

	return e
}

func errorHandler(err error, c echo.Context) {
	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		report, ok := err.(*echo.HTTPError)

		if !ok {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		switch report.Code {
		case http.StatusNotFound:
			response.ErrorBuilder(apperror.NotFound(errors.New("route not found"))).Send(c)

			return
		default:
			response.ErrorBuilder(err).Send(c)

			return
		}
	}

	// handle HTTP Error
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		response.ErrorBuilder(err).Send(c)
		return
	}
}

package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gruzowiki/rest/exceptions"
	"net/http"
)

type ErrorData struct {
	ErrorCode  string
	Message    string
	HttpStatus int
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"Message"`
}

func ErrorHandler(err error, c echo.Context) {
	var exception *exceptions.Exception
	if ok := errors.As(err, &exception); ok {
		errorData := handleMyException(*exception)
		c.JSON(errorData.HttpStatus, ErrorResponse{errorData.ErrorCode, errorData.Message})
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{"InternalServerError", "Internal Server Error"})
}

func handleMyException(exception exceptions.Exception) *ErrorData {
	if exception.Code == exceptions.IncorrectParams {
		return &ErrorData{
			exceptions.IncorrectParams,
			"Incorrect parameters provided: " + exception.Err.Error(),
			http.StatusBadRequest,
		}
	}

	if exception.Code == exceptions.CarrierNotFound {
		return &ErrorData{
			exceptions.CarrierNotFound,
			"Carrier not found: " + exception.Err.Error(),
			http.StatusNotFound,
		}
	}

	return &ErrorData{
		exceptions.InternalServerError,
		"Internal Server Error",
		http.StatusInternalServerError,
	}
}

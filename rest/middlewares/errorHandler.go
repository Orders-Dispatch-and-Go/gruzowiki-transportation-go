package middlewares

import (
	"errors"
	"gruzowiki/rest/terror"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err == nil {
			return nil
		}

		var code int

		var validationError terror.ValidationError
		var objectNotFound terror.NotFoundError
		switch {
		case errors.As(err, &validationError):
			code = http.StatusBadRequest
		case errors.As(err, &objectNotFound):
			code = http.StatusNotFound
		default:
			err = terror.NewInternalServerError(err.Error())
			code = http.StatusInternalServerError
		}

		return c.JSON(code, err)
	}
}

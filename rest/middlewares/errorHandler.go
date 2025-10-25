package middlewares

import (
	"gruzowiki/terror"
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

		switch err.(type) {
		case *echo.HTTPError:
			return err
		case terror.ValidationError:
			code = http.StatusBadRequest
		case terror.ObjectNotFound:
			code = http.StatusNotFound
		default:
			err = terror.NewInternalError(err.Error())
			code = http.StatusInternalServerError
		}

		return c.JSON(code, err)
	}
}
package middlewares

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gruzowiki/rest/terror"
	"io"
	"net/http"
	"strings"
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

		fmt.Println("Request ERROR")
		fmt.Printf("Method: %s URI: %s\n", c.Request().Method, c.Request().URL.Path)
		fmt.Printf("error: %s, code: %d\n", err, code)

		return c.JSON(code, err)
	}
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Printf("\nRequest Handle: %s %s", c.Request().Method, c.Request().URL.Path)
		fmt.Println("\nQuery:", c.Request().URL.RawQuery)
		//fmt.Println("Body:", c.Request().Body)

		/*if c.Request().Body != nil {
			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				fmt.Println("Error reading request body:", err)
			} else {
				c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				if len(bodyBytes) > 0 {
					fmt.Println("Body:")
					if json.Valid(bodyBytes) {
						var prettyJSON bytes.Buffer
						if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
							fmt.Println(prettyJSON.String())
						} else {
							fmt.Println(string(bodyBytes))
						}
					} else {
						fmt.Println(string(bodyBytes))
					}
				}
			}
		}*/

		if c.Request().Body != nil {
			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}
			oneLine := strings.ReplaceAll(string(bodyBytes), "\n", "")
			oneLine = strings.ReplaceAll(oneLine, " ", "")
			fmt.Println("Request body:", oneLine)
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		err := next(c)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return err
		}

		return nil
	}
}

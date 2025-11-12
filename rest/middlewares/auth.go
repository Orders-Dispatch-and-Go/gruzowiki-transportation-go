package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AllowedeRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format")
			}

			tokenString := parts[1]

			publicKeyData, err := os.ReadFile("public.key") //брать путь из конфига
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read public key")
			}

			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Invalid public key")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user_id in token")
			}

			role, ok := claims["role"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid role in token")
			}

			roleAllowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}

			c.Set("user_id", userID)
			c.Set("role", role)

			return next(c)
		}
	}
}
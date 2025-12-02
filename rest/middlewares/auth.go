package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	UserRole      = "ROLE_USER"
	ManagerRole   = "ROLE_MANAGER"
	AdminRole     = "ROLE_ADMIN"
	ConsignerRole = "ROLE_CONSIGNER"
	CarrierRole   = "ROLE_CARRIER"

	UserIdCtxClaim = "userId"
	EmailCtxClaim  = "email"
	RolesCtxClaim  = "roles"
)

type UserData struct {
	ID        int32  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"FirstName"`
}

func AllowedRoles(allowedRoles ...string) echo.MiddlewareFunc {
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

			publicKeyData, err := os.ReadFile("jws-public.pem") //брать путь из конфига
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

			/*
				roleAllowed := false
				for _, allowedRole := range allowedRoles {
					if role == allowedRole {
						roleAllowed = true
						break
					}
				}

				if !roleAllowed {
					return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
				} */

			userData, ok := claims["userData"].(map[string]interface{})
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid userData in token")
			}

			fmt.Println(userData["id"])

			rolesClaim, ok := claims["userAuthorities"].([]interface{})
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid userAuthorities in token")
			}

			roles := make([]string, 0)
			for _, value := range rolesClaim {
				roles = append(roles, value.(string))
			}

			c.Set(UserIdCtxClaim, int(userData["id"].(float64)))
			c.Set(EmailCtxClaim, userData["email"])
			c.Set(RolesCtxClaim, roles)

			return next(c)
		}
	}
}

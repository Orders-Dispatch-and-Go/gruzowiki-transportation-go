package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

type (
	JWTUserData struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
	}

	JWTPayload struct {
		Sub             string      `json:"sub"`
		UserAuthorities []string    `json:"userAuthorities"`
		UserData        JWTUserData `json:"userData"`
		IssuedAt        int64       `json:"issuedAt"`
		ExpiresAt       int64       `json:"expiresAt"`
		jwt.RegisteredClaims
	}
)

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

			/*publicKeyData, err := os.ReadFile("jws-public.pem") //брать путь из конфига
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

			userData, ok := claims["userData"].(map[string]interface{})
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid userData in token")
			}

			rolesClaim, ok := claims["userAuthorities"].([]interface{})
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid userAuthorities in token")
			}*/

			/*jwtParts := strings.Split(tokenString, ".")
			stringPayload := jwtParts[1]*/

			claims := &JWTPayload{}
			_, _, err := jwt.NewParser(jwt.WithoutClaimsValidation()).ParseUnverified(tokenString, claims)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}

			/*roles := make([]string, 0)
			for _, value := range rolesClaim {
				roles = append(roles, value.(string))
			}

			roleAllowed := false
			for _, role := range roles {
				if slices.Contains(allowedRoles, role) {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}*/

			fmt.Printf("UserJwtClaims: User ID: %d, Email: %s, Roles: %s", claims.UserData.ID, claims.UserData.Email, claims.UserAuthorities)
			c.Set(UserIdCtxClaim, claims.UserData.ID)
			c.Set(EmailCtxClaim, claims.UserData.Email)
			c.Set(RolesCtxClaim, claims.UserAuthorities)

			return next(c)
		}
	}
}

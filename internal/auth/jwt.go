package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
)

type Jwt struct {
	config JwtConfig
}

func NewJwt(config JwtConfig) Jwt {
	return Jwt{config: config}
}

func (j Jwt) ParseToken(tokenString string) (int32, error) {
	token, err := jwt.Parse(tokenString, j.getSecretKeyForToken)
	if err != nil {
		return 0, fmt.Errorf("parse: %w", err)
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("get user ID: %w", err)
	}

	userID, err := strconv.ParseInt(userIDString, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parse ID string %q: %w", userIDString, err)
	}

	return int32(userID), nil
}

func (j *Jwt) getSecretKeyForToken(token *jwt.Token) (any, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(j.config.publicKey), nil
}

type JwtConfig struct {
	publicKey string
}

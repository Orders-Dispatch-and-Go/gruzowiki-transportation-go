package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"auth-service/internal/utils/duration"
)

type JWT struct {
	cfg JWTConfig
}

func NewJWT(cfg JWTConfig) JWT {
	return JWT{cfg: cfg}
}

func (j JWT) GenerateToken(userID int64) (string, error) {
	userIDString := strconv.FormatInt(userID, 10)
	expiresAt := time.Now().Add(j.cfg.TokenLifetimeSeconds.Duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userIDString,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	tokenString, err := token.SignedString([]byte(j.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j JWT) GenerateVerificationToken(userID int64) (string, error) {
	userIDString := strconv.FormatInt(userID, 10)
	expiresAt := time.Now().Add(j.cfg.VerificationTokenLifetimeSeconds.Duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userIDString,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	tokenString, err := token.SignedString([]byte(j.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j JWT) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, j.getSecretKeyForToken)
	if err != nil {
		return 0, fmt.Errorf("parse: %w", err)
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("get user ID: %w", err)
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse ID string %q: %w", userIDString, err)
	}

	return userID, nil
}

func (j *JWT) getSecretKeyForToken(token *jwt.Token) (any, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(j.cfg.SecretKey), nil
}

type JWTConfig struct {
	SecretKey                        string           `validate:"required"`
	TokenLifetimeSeconds             duration.Seconds `validate:"required"`
	VerificationTokenLifetimeSeconds duration.Seconds `validate:"required"`
}

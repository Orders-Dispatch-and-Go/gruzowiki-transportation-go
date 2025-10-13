package middlewares

import "gruzowiki-transportation/internal/auth"

type Middlewares struct {
	auth auth.Jwt
}

func New(auth auth.Jwt) *Middlewares {
	return &Middlewares{auth: auth}
}

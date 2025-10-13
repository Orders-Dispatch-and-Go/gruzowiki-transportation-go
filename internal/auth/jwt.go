package auth

type Jwt struct {
	config JwtConfig
}

func NewJwt(config JwtConfig) Jwt {
	return Jwt{config: config}
}

func (j Jwt) ParseToken(token string) (int32, error) {
	return 0, nil
}

type JwtConfig struct {
	publicKey string
}

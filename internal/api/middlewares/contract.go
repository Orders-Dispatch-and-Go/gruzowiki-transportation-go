package middlewares

type auth interface {
	ParseToken(token string) (userID int64, err error)
}

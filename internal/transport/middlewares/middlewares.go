package middlewares

type Middlewares struct {
	auth auth
}

func New(auth auth) *Middlewares {
	return &Middlewares{auth: auth}
}

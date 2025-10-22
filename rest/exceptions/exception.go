package exceptions

type Exception struct {
	Code    string
	Message string
	Err     error
}

func (e *Exception) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewException(code string, err error) *Exception {
	return &Exception{
		Code: code,
		Err:  err,
	}
}

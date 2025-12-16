package terror

import "fmt"

type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func (e ErrorResponse) getType() string {
	return e.Type
}

type NotFoundError struct {
	ErrorResponse
	string
}

type ValidationError struct {
	ErrorResponse
}

type InternalServerError struct {
	ErrorResponse
}

func NewNotFoundError(objectType string, findByField string) NotFoundError {
	msg := fmt.Sprintf("Object of type %s not found by %s", objectType, findByField)
	return NotFoundError{
		ErrorResponse: ErrorResponse{
			Type:    TypeNotFoundError,
			Message: msg,
		},
	}
}

func NewValidationError(message string, cause string) ValidationError {
	msg := fmt.Sprintf("validation error: %s %s", cause, message)
	return ValidationError{
		ErrorResponse{
			Type:    TypeValidationError,
			Message: msg,
		},
	}
}

func NewInternalServerError(cause string) InternalServerError {
	msg := fmt.Sprintf("internal error: %s", cause)
	return InternalServerError{
		ErrorResponse{
			Type:    TypeInternalServerError,
			Message: msg,
		},
	}
}

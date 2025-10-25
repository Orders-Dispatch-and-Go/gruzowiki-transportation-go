package terror

/*
ФАЙЛ НЕ МЕНЯЕМ
*/

import (
 "encoding/json"
 "fmt"
 "reflect"
)

const (
	TypeInternalError   = "TerrorInternalError"
	TypeValidationError = "TerrorValidationError"
)

/*
ФАЙЛ НЕ МЕНЯЕМ
*/

type Base struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e Base) Error() string {
	return e.Message
}

func (e Base) gettype() string {
	return e.Type
}

type InternalError struct {
	Base
}

type ValidationError struct {
	Base
}

/*
ФАЙЛ НЕ МЕНЯЕМ
*/

func NewInternalError(cause string) InternalError {
	msg := fmt.Sprintf("internal error: %s", cause)
	return InternalError{
		Base: Base{
			Type:    TypeInternalError,
			Message: msg,
		},
	}
}

/*
ФАЙЛ НЕ МЕНЯЕМ
*/

func NewValidationError(cause string) ValidationError {
	msg := fmt.Sprintf("validation error: %s", cause)
	return ValidationError{
		Base: Base{
			Type:    TypeValidationError,
			Message: msg,
		},
	}
}

var commonErrInstances = map[string]interface{}{
	TypeInternalError:   InternalError{},
	TypeValidationError: ValidationError{},
}

/*
ФАЙЛ НЕ МЕНЯЕМ
*/

func unmarshal(body []byte, errInstances map[string]interface{}) (result error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			result = NewInternalError("invalid error instance")
		}
		}()
	var base Base
	if err := json.Unmarshal(body, &base); err != nil {
		return NewInternalError(string(body))
	}
	
	var instance interface{}
	if _, ok := errInstances[base.Type]; ok {
		instance = errInstances[base.Type]
	} else if _, ok := commonErrInstances[base.Type]; ok {
		instance = commonErrInstances[base.Type]
	}
	if instance != nil {
		target := reflect.New(reflect.TypeOf(instance))
		if err := json.Unmarshal(body, target.Interface()); err != nil {
			return NewInternalError(err.Error())
		}
		err, ok := target.Elem().Interface().(error)
		if !ok {
			return NewInternalError("invalid error type")
		}
		return err
	}
	return NewInternalError(string(body))
}

/*
ФАЙЛ НЕ МЕНЯЕМ
*/
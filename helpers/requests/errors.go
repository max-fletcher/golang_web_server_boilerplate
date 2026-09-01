package requests

import (
	"fmt"
)

type ErrInvalidUUIDParam struct {
	Param string
	Err   error
}

func (e ErrInvalidUUIDParam) Error() string {
	return fmt.Sprintf("Invalid UUID for route param %v: %v", e.Param, e.Err.Error())
}

type ErrInvalidJSON struct {
	Err error
}

func (e ErrInvalidJSON) Error() string {
	return fmt.Sprintf("Invalid JSON: %v", e.Err)
}

func (e ErrInvalidJSON) Unwrap() error {
	return e.Err
}

type ErrInvalidFieldType struct {
	Field string
	Type  string
}

func (e ErrInvalidFieldType) Error() string {
	return fmt.Sprintf(
		"%s must be of type %s",
		e.Field,
		e.Type,
	)
}

func (e ErrInvalidFieldType) FormatInvalidFieldTypeError() map[string]string {
	return map[string]string{
		e.Field: e.Error(),
	}
}

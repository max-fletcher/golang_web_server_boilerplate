package requests

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ErrInvalidJSON struct {
	Err error
}

func (e ErrInvalidJSON) Error() string {
	return fmt.Sprintf("invalid JSON: %v", e.Err)
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

func NewInvalidFieldTypeError(err error) error {
	var typeError *json.UnmarshalTypeError
	if !errors.As(err, &typeError) {
		return nil
	}

	return ErrInvalidFieldType{
		Field: typeError.Field,
		Type:  typeError.Type.String(),
	}
}

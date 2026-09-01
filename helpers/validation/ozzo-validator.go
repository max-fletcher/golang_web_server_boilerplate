// #TODO: Do I even need this
package validator

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func Validate(value validation.Validatable) error {
	err := value.Validate()
	if err == nil {
		return nil
	}

	return err
}

func FormatValidationErrors(err error) (map[string]string, bool) {
	if err == nil {
		return nil, false
	}

	// check if error matches type validation.Errors which is an interface for the errors that come from validating validatables
	var validationErrors validation.Errors
	if !errors.As(err, &validationErrors) {
		return nil, false
	}

	formattedErrors := make(map[string]string)
	for field, err := range validationErrors {
		formattedErrors[field] = err.Error()
	}

	return formattedErrors, true
}

func ValidateAndFormatValidationErrors(value validation.Validatable) (map[string]string, bool) {
	err := Validate(value)
	if err == nil {
		return nil, false
	}

	formattedErrors, ok := FormatValidationErrors(err)
	if !ok {
		return nil, false
	}

	return formattedErrors, true
}

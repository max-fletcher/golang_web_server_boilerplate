// #TODO: Do I even need this
package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func Validate(value validation.Validatable) map[string]string {
	err := value.Validate()
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validation.Errors)
	if !ok {
		return nil
	}

	errors := make(map[string]string)
	for field, err := range validationErrors {
		errors[field] = err.Error()
	}

	return errors
}

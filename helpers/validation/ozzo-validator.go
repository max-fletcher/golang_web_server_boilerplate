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
		return make(map[string]string), false
	}

	// check if error matches type validation.Errors which is an interface for the errors that come from validating validatables
	var validationErrors validation.Errors
	if !errors.As(err, &validationErrors) {
		return make(map[string]string), false
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

// custom validation rule used above
func PasswordsMatch(password string) validation.Rule {
	return validation.By(func(value interface{}) error {
		confirmPassword, ok := value.(string)
		if !ok {
			return errors.New("invalid password confirmation")
		}

		if confirmPassword != password {
			return errors.New("passwords do not match")
		}

		return nil
	})
}

func ConfirmPasswordRequiredOnPasswordProvided(password string) validation.Rule {
	return validation.When(password != "", validation.Required.Error("Confirm password is required"))
}
